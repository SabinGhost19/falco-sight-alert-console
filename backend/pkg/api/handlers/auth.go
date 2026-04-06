package handlers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login autentifică un utilizator și returnează un JWT token
func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Citim datele corecte din environment variables (fallback: admin/admin)
	expectedUser := os.Getenv("ADMIN_USER")
	if expectedUser == "" {
		expectedUser = "admin"
	}
	expectedPass := os.Getenv("ADMIN_PASSWORD")
	if expectedPass == "" {
		expectedPass = "admin"
	}

	// Verificăm credentialele
	if req.Username != expectedUser || req.Password != expectedPass {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	// Secretul JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-falcosight-key-change-in-production"
	}

	// Creăm claims-urile pentru token (valabil 24h)
	claims := jwt.MapClaims{
		"user":  req.Username,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	// Generăm token-ul
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{
		"token": t,
		"user":  req.Username,
	})
}
