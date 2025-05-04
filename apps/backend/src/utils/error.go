package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Status: code})
}

func GetValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s should be at least %s characters long", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s should be at most %s characters long", e.Field(), e.Param())
	default:
		return fmt.Sprintf("Failed validation on %s", e.Tag())
	}
}
