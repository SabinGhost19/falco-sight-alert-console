package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	customErrors "falcosight/pkg/api/errors"
	"falcosight/pkg/db"
	"falcosight/pkg/k8s"
	"falcosight/pkg/models"

	"github.com/gofiber/fiber/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HandleFalcoWebhook primește alertele de la Falco
func HandleFalcoWebhook(c *fiber.Ctx) error {
	var payload models.FalcoAlert
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Invalid Falco Alert Payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(customErrors.NewErrorResponse("ERR_INVALID_JSON", "Payload invalid de la Falco.", err.Error()))
	}

	// Extragem namespace și pod_name
	var namespace, podName, containerName string
	if val, ok := payload.OutputFields["k8s.ns.name"].(string); ok {
		namespace = val
	}
	if val, ok := payload.OutputFields["k8s.pod.name"].(string); ok {
		podName = val
	}
	if val, ok := payload.OutputFields["k8s.container.name"].(string); ok {
		containerName = val
	} else if val, ok := payload.OutputFields["container.name"].(string); ok {
		containerName = val
	}
	alertRecord := models.AlertModel{
		Priority:      payload.Priority,
		Rule:          payload.Rule,
		Message:       payload.Output,
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
	}

	// Dacă avem k8s pod/ns, pornim K8s Correlator (Analiză statică on-the-fly)
	if podName != "" && namespace != "" {
		if k8s.Clientset != nil {
			manifestYAML, err := k8s.FetchPodManifest(namespace, podName)
			if err != nil {
				log.Printf("Could not fetch manifest for %s/%s: %v", namespace, podName, err)
			} else {
				alertRecord.ManifestYAML = manifestYAML
				alertRecord.VulnerableLines = k8s.AnalyzeManifest(manifestYAML)
			}
		} else {
			log.Println("K8s Clientset is nil, skipping static analysis.")
		}
	}

	// Salvăm în baza de date
	if err := db.DB.Create(&alertRecord).Error; err != nil {
		log.Printf("Failed to save alert: %v", err)
		return customErrors.GlobalErrorHandler(c, err)
	}

	// ⏳ Forward Payload -> Falco Talon async
	talonURL := os.Getenv("TALON_WEBHOOK_URL")
	if talonURL == "" {
		talonURL = "http://falco-talon.falco-talon.svc.cluster.local:2803"
	}

	rawBody := c.Body()
	bodyCopy := make([]byte, len(rawBody))
	copy(bodyCopy, rawBody)

	go func(body []byte, url string) {
		client := &http.Client{Timeout: 10 * time.Second}
		httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/json")
			resp, reqErr := client.Do(httpReq)
			if reqErr == nil {
				defer resp.Body.Close()
				log.Printf("Forwarded alert to Talon at %s, status: %d", url, resp.StatusCode)
			} else {
				log.Printf("Failed forwarding alert to Talon: %v", reqErr)
			}
		}
	}(bodyCopy, talonURL)

	log.Printf("Ingested Falco Alert for Pod: %s", podName)
	return c.JSON(fiber.Map{"success": true, "id": alertRecord.ID})
}

// HandleTalonWebhook primește acțiunile executate de Falco Talon
func HandleTalonWebhook(c *fiber.Ctx) error {
	var payload models.TalonWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(customErrors.NewErrorResponse("ERR_INVALID_JSON", "Payload invalid de la Talon.", err.Error()))
	}

	var alert models.AlertModel
	res := db.DB.Where("pod_name = ? AND namespace = ?", payload.PodName, payload.Namespace).
		Order("created_at desc").First(&alert)

	if res.Error == nil {
		alert.TalonAction = payload.Action
		alert.TalonStatus = payload.Status
		db.DB.Save(&alert)
	}

	return c.JSON(fiber.Map{"success": true})
}

// GetAlerts funcție pentru Vue.js Frontend (Dashbaord Table)
func GetAlerts(c *fiber.Ctx) error {
	var alerts []models.AlertModel
	if err := db.DB.Order("created_at desc").Find(&alerts).Error; err != nil {
		return customErrors.GlobalErrorHandler(c, err)
	}
	return c.JSON(alerts)
}

// TriggerManualAction este apelată din interfață pentru a porni o remediere Talon la comandă
func TriggerManualAction(c *fiber.Ctx) error {
	var req struct {
		AlertID uint   `json:"alert_id"`
		Action  string `json:"action"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(customErrors.NewErrorResponse("ERR_VALIDATION", "Payload-ul acțiunii este invalid.", err.Error()))
	}

	var alert models.AlertModel
	if err := db.DB.First(&alert, req.AlertID).Error; err != nil {
		return customErrors.GlobalErrorHandler(c, err)
	}

	log.Printf("Triggering Talon Action: %s on Pod: %s for Alert ID: %d", req.Action, alert.PodName, alert.ID)

	// Validare via K8s (ca sa nu pornim talon pe un pod care deja nu mai exista si sa prindem RBAC errors on the fly)
	if k8s.Clientset != nil {
		_, err := k8s.Clientset.CoreV1().Pods(alert.Namespace).Get(context.Background(), alert.PodName, metav1.GetOptions{})
		if err != nil {
			return customErrors.GlobalErrorHandler(c, err)
		}
	}

	talonURL := os.Getenv("TALON_WEBHOOK_URL")
	if talonURL == "" {
		talonURL = "http://falco-talon.falco-talon.svc.cluster.local:2803"
	}

	payload := models.FalcoAlert{
		Output:   "Manual Override Triggered from FalcoSight Dashboard",
		Priority: "Critical",
		Rule:     "Manual Override: " + req.Action,
		Time:     time.Now().Format(time.RFC3339),
		OutputFields: map[string]interface{}{
			"k8s.pod.name":       alert.PodName,
			"k8s.container.name": alert.ContainerName,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return customErrors.GlobalErrorHandler(c, err)
	}

	// Folosim in mod obligatoriu contextul pentru a impune un hard-timeout de siguranta
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", talonURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return customErrors.GlobalErrorHandler(c, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		// Dacă Falco Talon este picat complet sau ajungem la timeout:
		return c.Status(fiber.StatusBadGateway).JSON(customErrors.NewErrorResponse(
			"ERR_TALON_GATEWAY",
			"Acțiunea nu poate fi executată temporar. Serviciul Falco Talon a picat sau nu a răspuns la timp.",
			err.Error(),
		))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return c.Status(fiber.StatusBadGateway).JSON(customErrors.NewErrorResponse(
			"ERR_TALON_INTERNAL",
			"Serviciul Falco Talon a returnat o eroare internă la procesarea comenzii.",
			fmt.Sprintf("Status Code: %d", resp.StatusCode),
		))
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		alert.TalonStatus = "Requested"
	} else {
		alert.TalonStatus = "Failed (Talon Error)"
	}

	alert.TalonAction = req.Action
	db.DB.Save(&alert)

	return c.JSON(fiber.Map{
		"success":      true,
		"pod":          alert.PodName,
		"action":       req.Action,
		"alertId":      alert.ID,
		"talon_status": alert.TalonStatus,
	})
}
