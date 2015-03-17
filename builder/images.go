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

func NewImageVars(img *models.Image, basePath string) *ImageVars {
	// eg: /site_1/image_m.jpg => /img/image_m.jpg
	return &ImageVars{
		Original:     path.Join(basePath, IMAGES_DIR, path.Base(img.URL())),
		Thumb:        path.Join(basePath, IMAGES_DIR, path.Base(img.ThumbURL())),
		Square:       path.Join(basePath, IMAGES_DIR, path.Base(img.SquareURL())),
		Small:        path.Join(basePath, IMAGES_DIR, path.Base(img.SmallURL())),
		SmallFill:    path.Join(basePath, IMAGES_DIR, path.Base(img.SmallFillURL())),
		PortraitFill: path.Join(basePath, IMAGES_DIR, path.Base(img.PortraitFillURL())),
		Large:        path.Join(basePath, IMAGES_DIR, path.Base(img.LargeURL())),
	}
}
