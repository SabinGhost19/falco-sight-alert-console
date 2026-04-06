package routes

import (
"os"

"falcosight/pkg/api/handlers"

jwtware "github.com/gofiber/contrib/jwt"
"github.com/gofiber/fiber/v2"
)

// Setup definește rutele aplicației
func Setup(app *fiber.App) {
	// --- Rute Publice (Ingest & Auth) ---
	// Ingest Webhooks (vin de la Falco/Talon - folosim autentificare tip MTLS sau Secret Token ideal, dar le lasam deschise pt simplitate)
	app.Post("/api/webhook/falco", handlers.HandleFalcoWebhook)
	app.Post("/api/webhook/talon", handlers.HandleTalonWebhook)

	// Login
	app.Post("/api/login", handlers.Login)

	// --- Rute Protejate cu JWT ---
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-falcosight-key-change-in-production"
	}

	// Grup protejat de API-uri pt Frontend
	v1 := app.Group("/api/v1")
	v1.Use(jwtware.New(jwtware.Config{
SigningKey: jwtware.SigningKey{Key: []byte(secret)},
ErrorHandler: func(c *fiber.Ctx, err error) error {
return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access, invalid or missing JWT"})
},
}))

	// API-uri consumate de UI (Vue.js) -> Necesită JWT
	v1.Get("/alerts", handlers.GetAlerts)
	v1.Post("/talon/trigger", handlers.TriggerManualAction)
}
