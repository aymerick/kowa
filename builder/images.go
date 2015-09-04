package builder

import (
	"path"

	"github.com/aymerick/kowa/models"
)

// ImageVars reprents an image variables
type ImageVars struct {
	Original             string
	OriginalAbsolute     string
	Thumb                string
	ThumbAbsolute        string
	Square               string
	SquareAbsolute       string
	Small                string
	SmallAbsolute        string
	SmallFill            string
	SmallFillAbsolute    string
	PortraitFill         string
	PortraitFillAbsolute string
	Large                string
	LargeAbsolute        string
}

// NewImageVars instanciates a new ImageVars
func NewImageVars(img *models.Image, basePath string, baseURL string) *ImageVars {
	// eg: image_m.jpg => /my_site/img/image_m.jpg
	//                    http://.../my_site/img/image_m.jpg
	return &ImageVars{
		Original:         path.Join("/", basePath, imagesDir, img.Path),
		OriginalAbsolute: baseURL + path.Join("/", imagesDir, img.Path),

		Thumb:         path.Join("/", basePath, imagesDir, img.ThumbPath()),
		ThumbAbsolute: baseURL + path.Join("/", imagesDir, img.ThumbPath()),

		Square:         path.Join("/", basePath, imagesDir, img.SquarePath()),
		SquareAbsolute: baseURL + path.Join("/", imagesDir, img.SquarePath()),

		Small:         path.Join("/", basePath, imagesDir, img.SmallPath()),
		SmallAbsolute: baseURL + path.Join("/", imagesDir, img.SmallPath()),

		SmallFill:         path.Join("/", basePath, imagesDir, img.SmallFillPath()),
		SmallFillAbsolute: baseURL + path.Join("/", imagesDir, img.SmallFillPath()),

		PortraitFill:         path.Join("/", basePath, imagesDir, img.PortraitFillPath()),
		PortraitFillAbsolute: baseURL + path.Join("/", imagesDir, img.PortraitFillPath()),

		Large:         path.Join("/", basePath, imagesDir, img.LargePath()),
		LargeAbsolute: baseURL + path.Join("/", imagesDir, img.LargePath()),
	}
}
