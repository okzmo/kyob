package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ParseAndValidate[T any](r *http.Request, validate *validator.Validate, body *T) error {
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return err
	}

	err = validate.Struct(body)
	if err != nil {
		return err
	}

	return nil
}
