package db

import (
	"log"

	"falcosight/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB inițializează conexiunea la baza de date
func ConnectDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Activează extensia în mod proactiv dacă baza de date e proaspătă
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb;`)

	// Auto-Migrate models
	err = DB.AutoMigrate(&models.AlertModel{}, &models.TalonRuleModel{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Converteste tabelul in Hypertable pentru stocare temporala masivă, partionat pe coloana `created_at`
	// Argumentul if_not_exists va ignora call-ul dacă este deja configurat
	DB.Exec(`SELECT create_hypertable('alert_models', 'created_at', if_not_exists => TRUE);`)

	log.Println("Database connection established, models migrated and TimescaleDB configured.")
}
