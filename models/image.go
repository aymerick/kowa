package models

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"time"

	"github.com/disintegration/imaging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	IMAGES_COL_NAME = "images"
	THUMB_SUFFIX    = "_t"
)

type Image struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`
	Path      string        `bson:"path"          json:"-"`
}

type ImagesList []Image

type ImageJson struct {
	Image
	URL      string `json:"url"`
	ThumbURL string `json:"thumbUrl"`
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
		Image:    *img,
		URL:      img.URL(),
		ThumbURL: img.ThumbURL(),
	}

	return json.Marshal(imageJson)
}

func (img *Image) ThumbPath() string {
	fileName := path.Base(img.Path)
	fileExt := path.Ext(img.Path)
	fileBase := fileName[:len(fileName)-len(fileExt)]

	return fmt.Sprintf("%s/%s%s%s", path.Dir(img.Path), fileBase, THUMB_SUFFIX, fileExt)
}

func (img *Image) URL() string {
	// @todo FIXME
	return img.Path
}

func (img *Image) ThumbURL() string {
	// @todo FIXME
	return img.ThumbPath()
}

func (img *Image) GenerateThumb(publicDir string) error {
	var err error
	var source image.Image

	sourcePath := path.Join(publicDir, img.Path)
	thumbPath := path.Join(publicDir, img.ThumbPath())

	// check if thumbnail already exists
	if _, err = os.Stat(thumbPath); !os.IsNotExist(err) {
		return nil
	}

	log.Printf("Generating thumb: %s", thumbPath)

	// open original
	source, err = imaging.Open(sourcePath)
	if err != nil {
		log.Printf("Failed to open: %v", sourcePath)
		return err
	}

	// create thumbnail
	thumb := imaging.Thumbnail(source, 100, 100, imaging.CatmullRom)

	// save thumbnail
	return imaging.Save(thumb, thumbPath)
}
