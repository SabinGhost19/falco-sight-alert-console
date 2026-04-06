package handlers

import (
"bytes"
"encoding/json"
"log"
"net/http"
"os"
"time"

"falcosight/pkg/db"
"falcosight/pkg/k8s"
"falcosight/pkg/models"

"github.com/gofiber/fiber/v2"
)

// HandleFalcoWebhook primește alertele de la Falco
func HandleFalcoWebhook(c *fiber.Ctx) error {
	var payload models.FalcoAlert
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("Invalid Falco Alert Payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Extragem namespace și pod_name
	namespace := payload.OutputFields["k8s.ns.name"]
	podName := payload.OutputFields["k8s.pod.name"]
	containerName := payload.OutputFields["k8s.container.name"]

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
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
	return c.JSON(fiber.Map{"status": "success", "id": alertRecord.ID})
}

// HandleTalonWebhook primește acțiunile executate de Falco Talon
func HandleTalonWebhook(c *fiber.Ctx) error {
	var payload models.TalonWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	// Căutăm alerta cea mai recentă asociată cu acest pod/namespace
	var alert models.AlertModel
	res := db.DB.Where("pod_name = ? AND namespace = ?", payload.PodName, payload.Namespace).
		Order("created_at desc").First(&alert)

	if res.Error == nil {
		// Actualizăm acțiunea Talon
		alert.TalonAction = payload.Action
		alert.TalonStatus = payload.Status
		db.DB.Save(&alert)
	}

	return c.JSON(fiber.Map{"status": "success"})
}

// GetAlerts funcție pentru Vue.js Frontend (Dashbaord Table)
func GetAlerts(c *fiber.Ctx) error {
	var alerts []models.AlertModel
	db.DB.Order("created_at desc").Find(&alerts)
	return c.JSON(alerts)
}

// TriggerManualAction este apelată din interfață pentru a porni o remediere Talon la comandă
func TriggerManualAction(c *fiber.Ctx) error {
	type RequestBody struct {
		AlertID uint   `json:"alert_id"`
		Action  string `json:"action"` // ex: "network_isolate", "terminate_pod"
	}

	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var alert models.AlertModel
	if err := db.DB.First(&alert, req.AlertID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alert not found"})
	}

	log.Printf("Triggering Talon Action: %s on Pod: %s for Alert ID: %d", req.Action, alert.PodName, alert.ID)

	talonURL := os.Getenv("TALON_WEBHOOK_URL")
	if talonURL == "" {
		talonURL = "http://falco-talon.falco-talon.svc.cluster.local:2803" // Default in cluster Kubernetes
	}

	payload := models.FalcoAlert{
		Output:   "Manual Override Triggered from FalcoSight Dashboard",
		Priority: "Critical",
		Rule:     "Manual Override: " + req.Action,
		Time:     time.Now().Format(time.RFC3339),
		OutputFields: map[string]string{
			"k8s.ns.name":        alert.Namespace,
			"k8s.pod.name":       alert.PodName,
			"k8s.container.name": alert.ContainerName,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal talon trigger payload: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request"})
	}

	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, err := http.NewRequest("POST", talonURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Failed to create HTTP request to Talon: %v", err)
	} else {
		httpReq.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(httpReq)
		if err != nil {
			log.Printf("Failed sending request to Falco Talon at %s: %v", talonURL, err)
			alert.TalonStatus = "Failed (Unreachable)"
		} else {
			defer resp.Body.Close()
			log.Printf("Sent request to Falco Talon, got status: %d", resp.StatusCode)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				alert.TalonStatus = "Requested"
			} else {
				alert.TalonStatus = "Failed (Talon Error)"
			}
		}
	}

	alert.TalonAction = req.Action
	db.DB.Save(&alert)

	return c.JSON(fiber.Map{
"status":  "Action Sent to Talon",
"pod":     alert.PodName,
"action":  req.Action,
"alertId": alert.ID,
"talon_status": alert.TalonStatus,
})
}
