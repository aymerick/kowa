package models

import (
	"encoding/json"
	"errors"
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

	"github.com/aymerick/kowa/utils"
)

const (
	IMAGES_COL_NAME = "images"

	// derivatives transformations kinds
	DERIVATIVE_FIT  = "fit"
	DERIVATIVE_FILL = "fill"

	// derivatives
	THUMB_KIND   = "thumb"
	THUMB_SCALE  = DERIVATIVE_FILL
	THUMB_SUFFIX = "_t"
	THUMB_WIDTH  = 100
	THUMB_HEIGHT = 75

	SQUARE_KIND   = "square"
	SQUARE_SCALE  = DERIVATIVE_FILL
	SQUARE_SUFFIX = "_q"
	SQUARE_WIDTH  = 200
	SQUARE_HEIGHT = 200

	SMALL_KIND   = "small"
	SMALL_SCALE  = DERIVATIVE_FIT
	SMALL_SUFFIX = "_s"
	SMALL_WIDTH  = 300
	SMALL_HEIGHT = 225

	SMALL_FILL_KIND   = "small_fill"
	SMALL_FILL_SCALE  = DERIVATIVE_FILL
	SMALL_FILL_SUFFIX = "_sf"
	SMALL_FILL_WIDTH  = 300
	SMALL_FILL_HEIGHT = 225

	PORTRAIT_FILL_KIND   = "portrait_fill"
	PORTRAIT_FILL_SCALE  = DERIVATIVE_FILL
	PORTRAIT_FILL_SUFFIX = "_pf"
	PORTRAIT_FILL_WIDTH  = 225
	PORTRAIT_FILL_HEIGHT = 300

	LARGE_KIND   = "large"
	LARGE_SCALE  = DERIVATIVE_FIT
	LARGE_SUFFIX = "_l"
	LARGE_WIDTH  = 1920
	LARGE_HEIGHT = 1440
)

type Image struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`
	Path      string        `bson:"path"          json:"-"`
	Name      string        `bson:"name"          json:"name"`
	Size      int64         `bson:"size"          json:"size"`
	Type      string        `bson:"type"          json:"type"`

	original *image.Image
}

type ImagesList []*Image

type ImageJson struct {
	Image
	URL string `json:"url"`

	ThumbURL        string `json:"thumbUrl"`
	SquareURL       string `json:"squareUrl"`
	SmallURL        string `json:"smallUrl"`
	SmallFillURL    string `json:"smallFillUrl"`
	PortraitFillURL string `json:"portraitFillUrl"`
	LargeURL        string `json:"largeUrl"`
}

type Derivative struct {
	kind   string
	scale  string
	suffix string
	width  int
	height int
}

var Derivatives []*Derivative

func init() {
	Derivatives = []*Derivative{
		&Derivative{
			kind:   THUMB_KIND,
			scale:  THUMB_SCALE,
			suffix: THUMB_SUFFIX,
			width:  THUMB_WIDTH,
			height: THUMB_HEIGHT,
		},
		&Derivative{
			kind:   SQUARE_KIND,
			scale:  SQUARE_SCALE,
			suffix: SQUARE_SUFFIX,
			width:  SQUARE_WIDTH,
			height: SQUARE_HEIGHT,
		},
		&Derivative{
			kind:   SMALL_KIND,
			scale:  SMALL_SCALE,
			suffix: SMALL_SUFFIX,
			width:  SMALL_WIDTH,
			height: SMALL_HEIGHT,
		},
		&Derivative{
			kind:   SMALL_FILL_KIND,
			scale:  SMALL_FILL_SCALE,
			suffix: SMALL_FILL_SUFFIX,
			width:  SMALL_FILL_WIDTH,
			height: SMALL_FILL_HEIGHT,
		},
		&Derivative{
			kind:   PORTRAIT_FILL_KIND,
			scale:  PORTRAIT_FILL_SCALE,
			suffix: PORTRAIT_FILL_SUFFIX,
			width:  PORTRAIT_FILL_WIDTH,
			height: PORTRAIT_FILL_HEIGHT,
		},
		&Derivative{
			kind:   LARGE_KIND,
			scale:  LARGE_SCALE,
			suffix: LARGE_SUFFIX,
			width:  LARGE_WIDTH,
			height: LARGE_HEIGHT,
		},
	}
}

// return a derivative definition
func DerivativeForKind(kind string) *Derivative {
	for _, derivative := range Derivatives {
		if derivative.kind == kind {
			return derivative
		}
	}

	return nil
}

// returns true if given path is an image derivative
func IsDerivativePath(path string) bool {
	for _, derivative := range Derivatives {
		fileBase := utils.FileBase(path)
		if strings.HasSuffix(fileBase, derivative.suffix) {
			return true
		}
	}

	return false
}

//
// DBSession
//

// Images collection
func (session *DBSession) ImagesCol() *mgo.Collection {
	return session.DB().C(IMAGES_COL_NAME)
}

// Ensure indexes on Images collection
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

// Find image by id
func (session *DBSession) FindImage(imageId bson.ObjectId) *Image {
	var result Image

	if err := session.ImagesCol().FindId(imageId).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// Persists a new image in database
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

// Implements json.MarshalJSON
func (img *Image) MarshalJSON() ([]byte, error) {

	imageJson := ImageJson{
		Image: *img,
		URL:   img.URL(),

		ThumbURL:        img.ThumbURL(),
		SquareURL:       img.SquareURL(),
		SmallURL:        img.SmallURL(),
		SmallFillURL:    img.SmallFillURL(),
		PortraitFillURL: img.PortraitFillURL(),
		LargeURL:        img.LargeURL(),
	}

	return json.Marshal(imageJson)
}

func (img *Image) Delete() error {
	var err error

	// delete from database
	if err = img.dbSession.ImagesCol().RemoveId(img.Id); err != nil {
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
	if err = os.Remove(originalPath); err != nil {
		log.Printf("Failed to delete image: %s", originalPath)
	}

	return nil
}

// Fetch from database: site that image belongs to
func (img *Image) FindSite() *Site {
	return img.dbSession.FindSite(img.SiteId)
}

// Returns memoized image original image
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

func (img *Image) OriginalFilePath() string {
	return path.Join(utils.AppPublicDir(), img.Path)
}

// Returns image URL
func (img *Image) URL() string {
	// @todo FIXME
	return img.Path
}

//
// Derivatives
//

// Returns Thumb derivative URL
func (img *Image) ThumbURL() string {
	return img.DerivativeURL(DerivativeForKind(THUMB_KIND))
}

// Returns Square derivative URL
func (img *Image) SquareURL() string {
	return img.DerivativeURL(DerivativeForKind(SQUARE_KIND))
}

// Returns Small derivative URL
func (img *Image) SmallURL() string {
	return img.DerivativeURL(DerivativeForKind(SMALL_KIND))
}

// Returns Small Fill derivative URL
func (img *Image) SmallFillURL() string {
	return img.DerivativeURL(DerivativeForKind(SMALL_FILL_KIND))
}

// Returns Portrait Fill derivative URL
func (img *Image) PortraitFillURL() string {
	return img.DerivativeURL(DerivativeForKind(PORTRAIT_FILL_KIND))
}

// Returns Large derivative URL
func (img *Image) LargeURL() string {
	return img.DerivativeURL(DerivativeForKind(LARGE_KIND))
}

func (img *Image) derivativePath(derivative *Derivative) string {
	return fmt.Sprintf("%s/%s%s%s", path.Dir(img.Path), utils.FileBase(img.Path), derivative.suffix, path.Ext(img.Path))
}

func (img *Image) DerivativeURL(derivative *Derivative) string {
	// @todo FIXME
	return img.derivativePath(derivative)
}

func (img *Image) DerivativeFilePath(derivative *Derivative) string {
	return path.Join(utils.AppPublicDir(), img.derivativePath(derivative))
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
	case DERIVATIVE_FIT:
		result = imaging.Fit(*img.Original(), derivative.width, derivative.height, imaging.Lanczos)

	case DERIVATIVE_FILL:
		result = imaging.Thumbnail(*img.Original(), derivative.width, derivative.height, imaging.Lanczos)

	default:
		panic("Insupported derivative scale")
	}

	// save derivative
	return imaging.Save(result, derivativePath)
}

// Generate all derivatives that were not generated yet
func (img *Image) GenerateDerivatives(force bool) error {
	var err error

	if img.Original() == nil {
		return errors.New(fmt.Sprintf("Failed to load original image: %v", img.Path))
	}

	for _, derivative := range Derivatives {
		if errGen := img.generateDerivative(derivative, force); errGen != nil {
			err = errGen
		}
	}

	return err
}
