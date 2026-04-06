package handlers

import (
	"falcosight/pkg/db"
	"falcosight/pkg/models"
	"github.com/gofiber/fiber/v2"
)

// GetRules => GET /api/v1/rules
func GetRules(c *fiber.Ctx) error {
	var rules []models.TalonRuleModel
	db.DB.Order("created_at desc").Find(&rules)
	return c.JSON(rules)
}

// CreateRule => POST /api/v1/rules
func CreateRule(c *fiber.Ctx) error {
	var rule models.TalonRuleModel
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	if err := db.DB.Create(&rule).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save rule"})
	}
	return c.JSON(rule)
}

// DeleteRule => DELETE /api/v1/rules/:id
func DeleteRule(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DB.Delete(&models.TalonRuleModel{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete rule"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Rule removed"})
}

// ToggleRule => PATCH /api/v1/rules/:id/toggle
func ToggleRule(c *fiber.Ctx) error {
	id := c.Params("id")
	var rule models.TalonRuleModel
	if err := db.DB.First(&rule, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}
	rule.Enabled = !rule.Enabled
	db.DB.Save(&rule)
	return c.JSON(rule)
}
