package models

import (
	"time"

	"gorm.io/gorm"
)

// FalcoAlert reprezintă structura payload-ului JSON venit prin webhook de la Falco.
type FalcoAlert struct {
	Output       string            `json:"output"`
	Priority     string            `json:"priority"`
	Rule         string            `json:"rule"`
	Time         string            `json:"time"`
	OutputFields map[string]string `json:"output_fields"`
}

// AlertModel este entitatea stocată în baza de date.
// Folosim direct structura pentru GORM, iar pentru TimescaleDB vom transforma tabelul in Hypertable dupa auto-migrate
type AlertModel struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt     time.Time      `gorm:"primaryKey;type:timestamptz;default:now()" json:"created_at"` // Composite PK req for Timescale
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Priority      string         `json:"priority"`
	Rule          string         `json:"rule"`
	Message       string         `json:"message"`
	Namespace     string         `json:"namespace"`
	PodName       string         `json:"pod_name"`
	ContainerName string         `json:"container_name"`

	// Correlated Data
	ManifestYAML    string `json:"manifest_yaml"`
	VulnerableLines string `json:"vulnerable_lines"` // Stochează indiciile de linii vulnerabile

	// Talon Correlation
	TalonAction string `json:"talon_action"`
	TalonStatus string `json:"talon_status"` // e.g., "Pending", "Success", "Failed"
}

// TalonWebhookPayload reprezintă payload-ul așteptat de la Falco Talon
type TalonWebhookPayload struct {
	Event     string `json:"event"`
	Action    string `json:"action"`
	PodName   string `json:"pod_name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
}
