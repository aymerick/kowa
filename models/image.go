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

	// derivatives
	SMALL_KIND   = "small"
	SMALL_SUFFIX = "_s"
	SMALL_WIDTH  = 100
	SMALL_HEIGHT = 60

	THUMB_KIND   = "thumb"
	THUMB_SUFFIX = "_t"
	THUMB_WIDTH  = 100
	THUMB_HEIGHT = 100

	MEDIUM_KIND   = "medium"
	MEDIUM_SUFFIX = "_m"
	MEDIUM_WIDTH  = 300
	MEDIUM_HEIGHT = 225

	MEDIUM_CROP_KIND   = "medium_crop"
	MEDIUM_CROP_SUFFIX = "_mc"
	MEDIUM_CROP_WIDTH  = 300
	MEDIUM_CROP_HEIGHT = 225
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
	URL           string `json:"url"`
	SmallURL      string `json:"smallUrl"`
	ThumbURL      string `json:"thumbUrl"`
	MediumURL     string `json:"mediumUrl"`
	MediumCropURL string `json:"mediumCropUrl"`
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
			kind:    SMALL_KIND,
			suffix:  SMALL_SUFFIX,
			genFunc: genSmall,
		},
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
		&Derivative{
			kind:    MEDIUM_CROP_KIND,
			suffix:  MEDIUM_CROP_SUFFIX,
			genFunc: genMediumCrop,
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
		Image:         *img,
		URL:           img.URL(),
		SmallURL:      img.derivativeURL(DerivativeForKind(SMALL_KIND)),
		ThumbURL:      img.derivativeURL(DerivativeForKind(THUMB_KIND)),
		MediumURL:     img.derivativeURL(DerivativeForKind(MEDIUM_KIND)),
		MediumCropURL: img.derivativeURL(DerivativeForKind(MEDIUM_CROP_KIND)),
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
		derivativePath := img.derivativeFilePath(derivative)

		if err := os.Remove(derivativePath); err != nil {
			log.Printf("Failed to delete image: %s", derivativePath)
		}
	}

	originalPath := img.originalFilePath()
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
		originalPath := img.originalFilePath()

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

func (img *Image) originalFilePath() string {
	return path.Join(appPublicDir, img.Path)
}

// Returns image URL
func (img *Image) URL() string {
	// @todo FIXME
	return img.Path
}

//
// Derivatives
//

func genSmall(source *image.Image) *image.NRGBA {
	return imaging.Thumbnail(*source, SMALL_WIDTH, SMALL_HEIGHT, imaging.Lanczos)
}

func genThumbnail(source *image.Image) *image.NRGBA {
	return imaging.Thumbnail(*source, THUMB_WIDTH, THUMB_HEIGHT, imaging.Lanczos)
}

func genMedium(source *image.Image) *image.NRGBA {
	return imaging.Fit(*source, MEDIUM_WIDTH, MEDIUM_HEIGHT, imaging.Lanczos)
}

func genMediumCrop(source *image.Image) *image.NRGBA {
	return imaging.Thumbnail(*source, MEDIUM_CROP_WIDTH, MEDIUM_CROP_HEIGHT, imaging.Lanczos)
}

func (img *Image) derivativePath(derivative *Derivative) string {
	return fmt.Sprintf("%s/%s%s%s", path.Dir(img.Path), utils.FileBase(img.Path), derivative.suffix, path.Ext(img.Path))
}

func (img *Image) derivativeURL(derivative *Derivative) string {
	// @todo FIXME
	return img.derivativePath(derivative)
}

func (img *Image) derivativeFilePath(derivative *Derivative) string {
	return path.Join(appPublicDir, img.derivativePath(derivative))
}

func (img *Image) generateDerivative(derivative *Derivative) error {
	derivativePath := img.derivativeFilePath(derivative)

	// check if derivative already exists
	if _, err := os.Stat(derivativePath); !os.IsNotExist(err) {
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
