package main

import (
"log"
"os"

"falcosight/pkg/api/routes"
"falcosight/pkg/db"
"falcosight/pkg/k8s"

"github.com/gofiber/fiber/v2"
"github.com/gofiber/fiber/v2/middleware/cors"
"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Inițializează K8s Client (Kubeconfig sau In-Cluster)
	k8s.InitK8sClient("")

	// 2. Conectarea la Baza de Date PostgreSQL. Citim dintr-un Env Var, alfel default value.
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "falcosight")
	dbPort := getEnv("DB_PORT", "5432")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
	db.ConnectDB(dsn)

	// 3. Inițializează aplicația Fiber
	app := fiber.New(fiber.Config{
AppName: "FalcoSight Ingestion API & Correlator",
})

	// Middleware-uri Globale
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
AllowOrigins: "*", // Pentru dezvoltare permitem tot.
AllowHeaders: "Origin, Content-Type, Accept",
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
