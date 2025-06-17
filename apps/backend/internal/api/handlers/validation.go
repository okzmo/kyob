package handlers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func SetupValidation() {
	validate = validator.New()
	validate.RegisterValidation("emoji_shortcode", validateEmojiShortcode)
}

func validateEmojiShortcode(fl validator.FieldLevel) bool {
	shortcode := fl.Field().String()

	// - Starts with [a-z]
	// - Ends with [a-z0-9]
	// - Middle can be [a-z0-9_] but no consecutive underscores
	// - Length 2-20
	pattern := `^[a-z]([a-z0-9]|_[a-z0-9])*[a-z0-9]$|^[a-z]$`

	return regexp.MustCompile(pattern).MatchString(shortcode)
}
