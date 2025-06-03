package utils

import (
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
