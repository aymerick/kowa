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
	// eg: /my_site/image_m.jpg => /my_site/img/image_m.jpg
	//                             http://.../my_site/img/image_m.jpg
	return &ImageVars{
		Original:             path.Join(basePath, IMAGES_DIR, path.Base(img.URL())),
		OriginalAbsolute:     baseUrl + path.Join("/", IMAGES_DIR, path.Base(img.URL())),
		Thumb:                path.Join(basePath, IMAGES_DIR, path.Base(img.ThumbURL())),
		ThumbAbsolute:        baseUrl + path.Join("/", IMAGES_DIR, path.Base(img.ThumbURL())),
		Square:               path.Join(basePath, IMAGES_DIR, path.Base(img.SquareURL())),
		SquareAbsolute:       baseUrl + path.Join("/", IMAGES_DIR, path.Base(img.SquareURL())),
		Small:                path.Join(basePath, IMAGES_DIR, path.Base(img.SmallURL())),
		SmallAbsolute:        baseUrl + path.Join("/", IMAGES_DIR, path.Base(img.SmallURL())),
		SmallFill:            path.Join(basePath, IMAGES_DIR, path.Base(img.SmallFillURL())),
		SmallFillAbsolute:    baseUrl + path.Join("/", IMAGES_DIR, path.Base(img.SmallFillURL())),
		PortraitFill:         path.Join(basePath, IMAGES_DIR, path.Base(img.PortraitFillURL())),
		PortraitFillAbsolute: baseUrl + path.Join("/", IMAGES_DIR, path.Base(img.PortraitFillURL())),
		Large:                path.Join(basePath, IMAGES_DIR, path.Base(img.LargeURL())),
		LargeAbsolute:        baseUrl + path.Join("/", IMAGES_DIR, path.Base(img.LargeURL())),
	}
}
