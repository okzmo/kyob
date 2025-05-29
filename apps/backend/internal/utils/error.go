package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	Code   string `json:"code"`
}

func RespondWithError(w http.ResponseWriter, status int, message string, code ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(code) <= 0 {
		json.NewEncoder(w).Encode(ErrorResponse{Error: message, Status: status})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Error: message, Status: status, Code: code[0]})
	}
}
