package utils

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/davidbyttow/govips/v2/vips"
)

func CropImage(file []byte, x, y, width, height int) ([]byte, error) {
	intSet := vips.IntParameter{}
	intSet.Set(-1)

	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(file, params)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	err = image.ExtractArea(x, y, width, height)
	if err != nil {
		return nil, err
	}

	webp := vips.NewWebpExportParams()
	webp.Lossless = false
	webp.NearLossless = false
	webp.Quality = 85
	webp.StripMetadata = true

	buf, _, err := image.ExportWebp(webp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ConvertToEmoji(file multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	intSet := vips.IntParameter{}
	intSet.Set(-1)

	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(fileBytes, params)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	width := image.Width()
	height := image.PageHeight()

	var scale float64
	if width > height {
		scale = float64(128) / float64(width)
	} else {
		scale = float64(128) / float64(height)
	}

	err = image.Resize(scale, vips.KernelLanczos3)
	if err != nil {
		return nil, err
	}

	webp := vips.NewWebpExportParams()
	webp.Lossless = false
	webp.NearLossless = false
	webp.Quality = 85
	webp.StripMetadata = true

	buf, _, err := image.ExportWebp(webp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ConvertToWebp(file multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	intSet := vips.IntParameter{}
	intSet.Set(-1)

	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(fileBytes, params)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	webp := vips.NewWebpExportParams()
	webp.Lossless = false
	webp.NearLossless = false
	webp.Quality = 85
	webp.StripMetadata = true

	buf, _, err := image.ExportWebp(webp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
