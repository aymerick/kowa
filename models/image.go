package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"time"

	"github.com/disintegration/imaging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/utils"
)

const (
	IMAGES_COL_NAME = "images"

	// derivatives
	THUMB_KIND    = "thumb"
	THUMB_SUFFIX  = "_t"
	MEDIUM_KIND   = "medium"
	MEDIUM_SUFFIX = "_m"
)

type Image struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`
	Path      string        `bson:"path"          json:"-"`

	original *image.Image
}

type ImagesList []Image

type ImageJson struct {
	Image
	URL       string `json:"url"`
	ThumbURL  string `json:"thumbUrl"`
	MediumURL string `json:"mediumUrl"`
}

type DerivativeGenFunc func(source *image.Image) *image.NRGBA

type Derivative struct {
	kind    string
	suffix  string
	genFunc DerivativeGenFunc
}

// @todo FIXME
var appPublicDir string

var Derivatives []*Derivative

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	appPublicDir = path.Join(currentDir, "/client/public")

	Derivatives = []*Derivative{
		&Derivative{
			kind:    THUMB_KIND,
			suffix:  THUMB_SUFFIX,
			genFunc: genThumbnail,
		},
		&Derivative{
			kind:    MEDIUM_KIND,
			suffix:  MEDIUM_SUFFIX,
			genFunc: genMedium,
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

//
// Image
//

// Implements json.MarshalJSON
func (img *Image) MarshalJSON() ([]byte, error) {

	imageJson := ImageJson{
		Image:     *img,
		URL:       img.URL(),
		ThumbURL:  img.derivativeURL(DerivativeForKind("thumb")),
		MediumURL: img.derivativeURL(DerivativeForKind("medium")),
	}

	return json.Marshal(imageJson)
}

// Returns memoized image original image
func (img *Image) Original() *image.Image {
	if img.original == nil {
		originalPath := path.Join(appPublicDir, img.Path)

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

// Returns image URL
func (img *Image) URL() string {
	// @todo FIXME
	return img.Path
}

//
// Derivatives
//

func genThumbnail(source *image.Image) *image.NRGBA {
	return imaging.Thumbnail(*source, 100, 100, imaging.Lanczos)
}

func genMedium(source *image.Image) *image.NRGBA {
	return imaging.Fit(*source, 200, 200, imaging.Lanczos)
}

func (img *Image) derivativePath(derivative *Derivative) string {
	return fmt.Sprintf("%s/%s%s%s", path.Dir(img.Path), utils.FileBase(img.Path), derivative.suffix, path.Ext(img.Path))
}

func (img *Image) derivativeURL(derivative *Derivative) string {
	// @todo FIXME
	return img.derivativePath(derivative)
}

func (img *Image) generateDerivative(derivative *Derivative) error {
	var err error

	derivativePath := path.Join(appPublicDir, img.derivativePath(derivative))

	// check if derivative already exists
	if _, err = os.Stat(derivativePath); !os.IsNotExist(err) {
		return nil
	}

	log.Printf("Generating derivative %s: %s", derivative.kind, derivativePath)

	// create derivative
	result := derivative.genFunc(img.Original())

	// save derivative
	return imaging.Save(result, derivativePath)
}

// Generate all derivatives that were not generated yet
func (img *Image) GenerateDerivatives() error {
	var err error

	if img.Original() == nil {
		return errors.New(fmt.Sprintf("Failed to load original image: %v", img.Path))
	}

	for _, derivative := range Derivatives {
		if errGen := img.generateDerivative(derivative); errGen != nil {
			err = errGen
		}
	}

	return err
}
