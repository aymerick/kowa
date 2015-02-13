package builder

import (
	"path"

	"github.com/aymerick/kowa/models"
)

// Image vars
type ImageVars struct {
	URL           string
	ThumbFillURL  string
	SquareFillURL string
	SmallURL      string
	SmallFillURL  string
	LargeURL      string
}

func NewImageVars(img *models.Image) *ImageVars {
	// eg: /site_1/image_m.jpg => /img/image_m.jpg
	return &ImageVars{
		URL:           path.Join("/", IMAGES_DIR, path.Base(img.URL())),
		ThumbFillURL:  path.Join("/", IMAGES_DIR, path.Base(img.ThumbFillURL())),
		SquareFillURL: path.Join("/", IMAGES_DIR, path.Base(img.SquareFillURL())),
		SmallURL:      path.Join("/", IMAGES_DIR, path.Base(img.SmallURL())),
		SmallFillURL:  path.Join("/", IMAGES_DIR, path.Base(img.SmallFillURL())),
		LargeURL:      path.Join("/", IMAGES_DIR, path.Base(img.LargeURL())),
	}
}
