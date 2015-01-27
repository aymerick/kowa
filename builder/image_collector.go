package builder

import (
	"fmt"

	"github.com/aymerick/kowa/models"
)

type ImageKind struct {
	Image *models.Image
	Kind  string
}

type ImageCollector struct {
	Images map[string]*ImageKind
}

func NewImageCollector() *ImageCollector {
	return &ImageCollector{
		Images: make(map[string]*ImageKind),
	}
}

func NewImageKind(img *models.Image, kind string) *ImageKind {
	return &ImageKind{
		Image: img,
		Kind:  kind,
	}
}

// Add a new error for given step
func (collector *ImageCollector) addImage(img *models.Image, kind string) {
	key := fmt.Sprintf("%s-%s", img.Id.String(), kind)
	collector.Images[key] = NewImageKind(img, kind)
}
