package utils

import (
	"mime/multipart"

	"github.com/davidbyttow/govips/v2/vips"
)

func CropImage(file multipart.File, x, y, width, height int) ([]byte, error) {
	image, err := vips.NewImageFromReader(file)
	if err != nil {
		return nil, err
	}

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
