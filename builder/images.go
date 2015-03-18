package builder

import (
	"path"

	"github.com/aymerick/kowa/models"
)

// Image vars
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

func NewImageVars(img *models.Image, basePath string, baseUrl string) *ImageVars {
	// eg: image_m.jpg => /my_site/img/image_m.jpg
	//                    http://.../my_site/img/image_m.jpg
	return &ImageVars{
		Original:         path.Join(basePath, IMAGES_DIR, img.Path),
		OriginalAbsolute: baseUrl + path.Join("/", IMAGES_DIR, img.Path),

		Thumb:         path.Join(basePath, IMAGES_DIR, img.ThumbPath()),
		ThumbAbsolute: baseUrl + path.Join("/", IMAGES_DIR, img.ThumbPath()),

		Square:         path.Join(basePath, IMAGES_DIR, img.SquarePath()),
		SquareAbsolute: baseUrl + path.Join("/", IMAGES_DIR, img.SquarePath()),

		Small:         path.Join(basePath, IMAGES_DIR, img.SmallPath()),
		SmallAbsolute: baseUrl + path.Join("/", IMAGES_DIR, img.SmallPath()),

		SmallFill:         path.Join(basePath, IMAGES_DIR, img.SmallFillPath()),
		SmallFillAbsolute: baseUrl + path.Join("/", IMAGES_DIR, img.SmallFillPath()),

		PortraitFill:         path.Join(basePath, IMAGES_DIR, img.PortraitFillPath()),
		PortraitFillAbsolute: baseUrl + path.Join("/", IMAGES_DIR, img.PortraitFillPath()),

		Large:         path.Join(basePath, IMAGES_DIR, img.LargePath()),
		LargeAbsolute: baseUrl + path.Join("/", IMAGES_DIR, img.LargePath()),
	}
}
