package errors

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// ErrorResponse contract standard
type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	} `json:"error"`
}

// NewErrorResponse returneaza o structura JSON pentru erori previzibile
func NewErrorResponse(code, message, details string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		}{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// GlobalErrorHandler prinde toate erorile din rute si le formateaza corespunzator
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	resp := NewErrorResponse("ERR_INTERNAL_SERVER", "A apărut o eroare internă neașteptată.", err.Error())

	// Verificăm erori K8s
	if k8serrors.IsNotFound(err) {
		code = fiber.StatusNotFound
		resp = NewErrorResponse("ERR_K8S_NOT_FOUND", "Resursa Kubernetes nu a fost găsită (posibil ștearsă).", err.Error())
	} else if k8serrors.IsForbidden(err) {
		code = fiber.StatusForbidden
		resp = NewErrorResponse("ERR_K8S_FORBIDDEN", "Aplicația nu are permisiuni RBAC pentru a contacta Kubernetes API.", err.Error())
	} else if k8serrors.IsTimeout(err) {
		code = fiber.StatusGatewayTimeout
		resp = NewErrorResponse("ERR_K8S_TIMEOUT", "Kubernetes API a returnat un Timeout.", err.Error())
	} else if err == gorm.ErrRecordNotFound {
		// Erori GORM baza de date
		code = fiber.StatusNotFound
		resp = NewErrorResponse("ERR_DB_NOT_FOUND", "Înregistrarea solicitată nu există în sistem.", err.Error())
	}

	// Alte erori Fiber standard
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		resp = NewErrorResponse("ERR_HTTP", e.Message, "")
	}

	return c.Status(code).JSON(resp)
}
