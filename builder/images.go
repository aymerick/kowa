package builder

import (
	"path"

	"github.com/aymerick/kowa/models"
)

// Image vars
type ImageVars struct {
	Original     string
	Thumb        string
	Square       string
	Small        string
	SmallFill    string
	PortraitFill string
	Large        string
}

func NewImageVars(img *models.Image, baseURL string) *ImageVars {
	// eg: /site_1/image_m.jpg => /img/image_m.jpg
	return &ImageVars{
		Original:     path.Join(baseURL, IMAGES_DIR, path.Base(img.URL())),
		Thumb:        path.Join(baseURL, IMAGES_DIR, path.Base(img.ThumbURL())),
		Square:       path.Join(baseURL, IMAGES_DIR, path.Base(img.SquareURL())),
		Small:        path.Join(baseURL, IMAGES_DIR, path.Base(img.SmallURL())),
		SmallFill:    path.Join(baseURL, IMAGES_DIR, path.Base(img.SmallFillURL())),
		PortraitFill: path.Join(baseURL, IMAGES_DIR, path.Base(img.PortraitFillURL())),
		Large:        path.Join(baseURL, IMAGES_DIR, path.Base(img.LargeURL())),
	}
}
