package utils

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/go-playground/validator/v10"
)

type ImageValidationConfig struct {
	MaxSize            int64
	AllowedMimeTypes   []string
	RequireValidHeader bool
}

func ParseAndValidateImage(fileHeader *multipart.FileHeader, config ImageValidationConfig) error {
	if fileHeader.Size > config.MaxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", fileHeader.Size, config.MaxSize)
	}

	if fileHeader.Size == 0 {
		return fmt.Errorf("empty file not allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	detectedMimeType := http.DetectContentType(buffer[:n])

	if !isAllowedMimeType(detectedMimeType, config.AllowedMimeTypes) {
		return fmt.Errorf("unsupported MIME type: %s", detectedMimeType)
	}

	if seeker, ok := file.(io.Seeker); ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to reset file pointer: %w", err)
		}
	}

	if config.RequireValidHeader {
		if err := validateImageStructure(file); err != nil {
			return fmt.Errorf("invalid image structure: %w", err)
		}
	}

	return nil
}

func isAllowedMimeType(mimeType string, allowed []string) bool {
	return slices.Contains(allowed, mimeType)
}

func validateImageStructure(file multipart.File) error {
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	_, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	switch format {
	case "jpeg", "png", "gif", "webp":
		break
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}

	return nil
}

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
