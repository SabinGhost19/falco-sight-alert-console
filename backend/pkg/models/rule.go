package models

import (
	"gorm.io/gorm"
)

// TalonRuleModel este entitatea pentru stocarea regulilor (mappings: Falco -> Talon Action)
type TalonRuleModel struct {
	gorm.Model
	Name          string `json:"name"`          // ex: "Isolate Miner Pod"
	FalcoRule     string `json:"falcoRule"`     // ex: "Terminal shell in container"
	Action        string `json:"action"`        // ex: "network_isolate", "terminate_pod", "label_pod"
	ActionDetails string `json:"actionDetails"` // ex: "label: quarantined=true" (opțional)
	Enabled       bool   `json:"enabled"`
}
