package main

import (
	"log"
	"os"

	"falcosight/pkg/api/errors"
	"falcosight/pkg/api/routes"
	"falcosight/pkg/db"
	"falcosight/pkg/k8s"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Inițializează K8s Client (Kubeconfig sau In-Cluster)
	k8s.InitK8sClient(os.Getenv("KUBECONFIG"))

	// 2. Conectarea la Baza de Date PostgreSQL. Citim dintr-un Env Var, alfel default value.
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "falcosight")
	dbPort := getEnv("DB_PORT", "5432")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
	db.ConnectDB(dsn)

	// 3. Inițializează aplicația Fiber + Global Error Handler
	app := fiber.New(fiber.Config{
		AppName:      "FalcoSight Ingestion API & Correlator",
		ErrorHandler: errors.GlobalErrorHandler,
	})

	// Verificari de Securitate Production Ready
	if os.Getenv("JWT_SECRET") == "" || os.Getenv("ADMIN_USER") == "" || os.Getenv("ADMIN_PASSWORD") == "" {
		log.Fatal("CRITICAL: Missing JWT_SECRET, ADMIN_USER, or ADMIN_PASSWORD. Production environment must provide secure credentials.")
	}

	// Middleware-uri Globale (Recover is critical for enterprise)
	app.Use(recover.New())
	app.Use(logger.New())

	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Println("WARN: CORS_ALLOWED_ORIGINS not set. Defaulting to strict same-origin only.")
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 4. Setarea Rutelor (Webhook-uri Falco + API UI)
	routes.Setup(app)

	// 5. Pornirea Serverului
	port := getEnv("PORT", "3000")
	log.Printf("Starting FalcoSight Backend on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// getEnv citeste o valoare din structura de environment. Daca nu exista trimite fallback-ul.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
