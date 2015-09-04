package models

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/helpers"
)

const (
	imagesColName = "images"

	// derivatives transformations kinds
	derivativeFit  = "fit"
	derivativeFill = "fill"

	// derivatives
	thumbKind   = "thumb"
	thumbScale  = derivativeFill
	thumbSuffix = "_t"
	thumbWidth  = 100
	thumbHeight = 75

	squareKind   = "square"
	squareScale  = derivativeFill
	squareSuffix = "_q"
	squareWidth  = 200
	squareHeight = 200

	smallKind   = "small"
	smallScale  = derivativeFit
	smallSuffix = "_s"
	smallWidth  = 300
	smallHeight = 225

	smallFillKind   = "small_fill"
	smallFillScale  = derivativeFill
	smallFillSuffix = "_sf"
	smallFillWidth  = 300
	smallFillHeight = 225

	portraitFillKind   = "portrait_fill"
	portraitFillScale  = derivativeFill
	portraitFillSuffix = "_pf"
	portraitFillWidth  = 225
	portraitFillHeight = 300

	largeKind   = "large"
	largeScale  = derivativeFit
	largeSuffix = "_l"
	largeWidth  = 1920
	largeHeight = 1440
)

// Image represents an image
type Image struct {
	dbSession *DBSession `bson:"-"`

	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteID    string        `bson:"site_id"       json:"site"`
	Path      string        `bson:"path"          json:"-"`    // this is the effective image path
	Name      string        `bson:"name"          json:"name"` // this is the uploaded file name (may be different from Path)
	Size      int64         `bson:"size"          json:"size"`
	Type      string        `bson:"type"          json:"type"` // jpeg | png

	original *image.Image
}

// ImagesList represents an list of images
type ImagesList []*Image

// ImageJSON represents the json version of an image
type ImageJSON struct {
	Image
	URL string `json:"url"`

	ThumbURL        string `json:"thumbUrl"`
	SquareURL       string `json:"squareUrl"`
	SmallURL        string `json:"smallUrl"`
	SmallFillURL    string `json:"smallFillUrl"`
	PortraitFillURL string `json:"portraitFillUrl"`
	LargeURL        string `json:"largeUrl"`
}

// Derivative represents an image derivative
type Derivative struct {
	kind   string
	scale  string
	suffix string
	width  int
	height int
}

// Derivatives represents a list of image derivatives
var Derivatives []*Derivative

func init() {
	Derivatives = []*Derivative{
		&Derivative{
			kind:   thumbKind,
			scale:  thumbScale,
			suffix: thumbSuffix,
			width:  thumbWidth,
			height: thumbHeight,
		},
		&Derivative{
			kind:   squareKind,
			scale:  squareScale,
			suffix: squareSuffix,
			width:  squareWidth,
			height: squareHeight,
		},
		&Derivative{
			kind:   smallKind,
			scale:  smallScale,
			suffix: smallSuffix,
			width:  smallWidth,
			height: smallHeight,
		},
		&Derivative{
			kind:   smallFillKind,
			scale:  smallFillScale,
			suffix: smallFillSuffix,
			width:  smallFillWidth,
			height: smallFillHeight,
		},
		&Derivative{
			kind:   portraitFillKind,
			scale:  portraitFillScale,
			suffix: portraitFillSuffix,
			width:  portraitFillWidth,
			height: portraitFillHeight,
		},
		&Derivative{
			kind:   largeKind,
			scale:  largeScale,
			suffix: largeSuffix,
			width:  largeWidth,
			height: largeHeight,
		},
	}
}

// DerivativeForKind returns a derivative definition
func DerivativeForKind(kind string) *Derivative {
	for _, derivative := range Derivatives {
		if derivative.kind == kind {
			return derivative
		}
	}

	return nil
}

// IsDerivativePath returns true if given path is an image derivative
func IsDerivativePath(path string) bool {
	for _, derivative := range Derivatives {
		fileBase := helpers.FileBase(path)
		if strings.HasSuffix(fileBase, derivative.suffix) {
			return true
		}
	}

	return false
}

//
// DBSession
//

// ImagesCol returns images collection
func (session *DBSession) ImagesCol() *mgo.Collection {
	return session.DB().C(imagesColName)
}

// EnsureImagesIndexes ensures indexes on images collection
func (session *DBSession) EnsureImagesIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.ImagesCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// FindImage finds an image by id
func (session *DBSession) FindImage(imageID bson.ObjectId) *Image {
	var result Image

	if err := session.ImagesCol().FindId(imageID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreateImage creates a new image in database
func (session *DBSession) CreateImage(img *Image) error {
	now := time.Now()
	img.CreatedAt = now
	img.UpdatedAt = now

	if err := session.ImagesCol().Insert(img); err != nil {
		return err
	}

	img.dbSession = session

	return nil
}

//
// Image
//

// MarshalJSON implements the json.Marshaler interface
func (img *Image) MarshalJSON() ([]byte, error) {
	imageJSON := ImageJSON{
		Image: *img,
		URL:   img.URL(),

		ThumbURL:        img.ThumbURL(),
		SquareURL:       img.SquareURL(),
		SmallURL:        img.SmallURL(),
		SmallFillURL:    img.SmallFillURL(),
		PortraitFillURL: img.PortraitFillURL(),
		LargeURL:        img.LargeURL(),
	}

	return json.Marshal(imageJSON)
}

// Delete deletes image from database
func (img *Image) Delete() error {
	// delete from database
	if err := img.dbSession.ImagesCol().RemoveId(img.ID); err != nil {
		return err
	}

	// delete files
	for _, derivative := range Derivatives {
		derivativePath := img.DerivativeFilePath(derivative)

		if err := os.Remove(derivativePath); err != nil {
			log.Printf("Failed to delete image: %s", derivativePath)
		}
	}

	originalPath := img.OriginalFilePath()
	if err := os.Remove(originalPath); err != nil {
		log.Printf("Failed to delete image: %s", originalPath)
	}

	return nil
}

// FindSite fetches site that image belongs to
func (img *Image) FindSite() *Site {
	return img.dbSession.FindSite(img.SiteID)
}

// Original returns memoized image original image
func (img *Image) Original() *image.Image {
	if img.original == nil {
		originalPath := img.OriginalFilePath()

		// open original
		openedImage, err := imaging.Open(originalPath)
		if err != nil {
			log.Printf("Failed to open: %v", originalPath)
		} else {
			img.original = &openedImage
		}
	}

	return img.original
}

// OriginalFilePath returns file path to original image
func (img *Image) OriginalFilePath() string {
	return core.UploadSiteFilePath(img.SiteID, img.Path)
}

// URL returns image URL
func (img *Image) URL() string {
	return core.UploadSiteUrlPath(img.SiteID, img.Path)
}

//
// Derivatives
//

// ThumbPath returns Thumb derivative path
func (img *Image) ThumbPath() string {
	return img.DerivativePath(DerivativeForKind(thumbKind))
}

// ThumbURL returns Thumb derivative URL
func (img *Image) ThumbURL() string {
	return img.DerivativeURL(DerivativeForKind(thumbKind))
}

// SquarePath returns Square derivative path
func (img *Image) SquarePath() string {
	return img.DerivativePath(DerivativeForKind(squareKind))
}

// SquareURL returns Square derivative URL
func (img *Image) SquareURL() string {
	return img.DerivativeURL(DerivativeForKind(squareKind))
}

// SmallPath returns Small derivative path
func (img *Image) SmallPath() string {
	return img.DerivativePath(DerivativeForKind(smallKind))
}

// SmallURL returns Small derivative URL
func (img *Image) SmallURL() string {
	return img.DerivativeURL(DerivativeForKind(smallKind))
}

// SmallFillPath returns SmallFill derivative path
func (img *Image) SmallFillPath() string {
	return img.DerivativePath(DerivativeForKind(smallFillKind))
}

// SmallFillURL returns Small Fill derivative URL
func (img *Image) SmallFillURL() string {
	return img.DerivativeURL(DerivativeForKind(smallFillKind))
}

// PortraitFillPath returns PortraitFill derivative path
func (img *Image) PortraitFillPath() string {
	return img.DerivativePath(DerivativeForKind(portraitFillKind))
}

// PortraitFillURL returns Portrait Fill derivative URL
func (img *Image) PortraitFillURL() string {
	return img.DerivativeURL(DerivativeForKind(portraitFillKind))
}

// LargePath returns Large derivative path
func (img *Image) LargePath() string {
	return img.DerivativePath(DerivativeForKind(largeKind))
}

// LargeURL returns Large derivative URL
func (img *Image) LargeURL() string {
	return img.DerivativeURL(DerivativeForKind(largeKind))
}

// DerivativePath returns given derivative path
func (img *Image) DerivativePath(derivative *Derivative) string {
	return fmt.Sprintf("%s%s%s", helpers.FileBase(img.Path), derivative.suffix, path.Ext(img.Path))
}

// DerivativeURL returns given derivative URL
func (img *Image) DerivativeURL(derivative *Derivative) string {
	return core.UploadSiteUrlPath(img.SiteID, img.DerivativePath(derivative))
}

// DerivativeFilePath returns given derivative file path
func (img *Image) DerivativeFilePath(derivative *Derivative) string {
	return core.UploadSiteFilePath(img.SiteID, img.DerivativePath(derivative))
}

func (img *Image) generateDerivative(derivative *Derivative, force bool) error {
	derivativePath := img.DerivativeFilePath(derivative)

	if !force {
		// check if derivative already exists
		if _, err := os.Stat(derivativePath); !os.IsNotExist(err) {
			return nil
		}
	}

	log.Printf("Generating derivative %s: %s", derivative.kind, derivativePath)

	// create derivative
	var result *image.NRGBA

	switch derivative.scale {
	case derivativeFit:
		result = imaging.Fit(*img.Original(), derivative.width, derivative.height, imaging.Lanczos)

	case derivativeFill:
		result = imaging.Thumbnail(*img.Original(), derivative.width, derivative.height, imaging.Lanczos)

	default:
		panic("Insupported derivative scale")
	}

	// save derivative
	return imaging.Save(result, derivativePath)
}

// GenerateDerivatives generates all derivatives that were not generated yet
func (img *Image) GenerateDerivatives(force bool) error {
	var err error

	if img.Original() == nil {
		return fmt.Errorf("Failed to load original image: %v", img.Path)
	}

	for _, derivative := range Derivatives {
		if errGen := img.generateDerivative(derivative, force); errGen != nil {
			err = errGen
		}
	}

	return err
}
