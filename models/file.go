package models

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aymerick/kowa/core"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	FILES_COL_NAME  = "files"
	FILE_MEMBERSHIP = "membership"
)

var FileKinds = []string{FILE_MEMBERSHIP}

type File struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	Kind string `bson:"kind" json:"kind"`
	Path string `bson:"path" json:"-"`    // this is the effective file path
	Name string `bson:"name" json:"name"` // this is the uploaded file name (may be different from Path)
	Size int64  `bson:"size" json:"size"`
	Type string `bson:"type" json:"type"` // content type
}

type FilesList []*File

type FileJson struct {
	File
	URL string `json:"url"`
}

func IsValidFileKind(kind string) bool {
	for _, k := range FileKinds {
		if k == kind {
			return true
		}
	}

	return false
}

//
// DBSession
//

// Files collection
func (session *DBSession) FilesCol() *mgo.Collection {
	return session.DB().C(FILES_COL_NAME)
}

// Ensure indexes on Files collection
func (session *DBSession) EnsureFilesIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.FilesCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// Find file by id
func (session *DBSession) FindFile(fileId bson.ObjectId) *File {
	var result File

	if err := session.FilesCol().FindId(fileId).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// Persists a new file in database
func (session *DBSession) CreateFile(f *File) error {
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now

	if err := session.FilesCol().Insert(f); err != nil {
		return err
	}

	f.dbSession = session

	return nil
}

//
// File
//

// Implements json.MarshalJSON
func (f *File) MarshalJSON() ([]byte, error) {

	fileJson := FileJson{
		File: *f,
		URL:  f.URL(),
	}

	return json.Marshal(fileJson)
}

func (f *File) Delete() error {
	// delete from database
	if err := f.dbSession.FilesCol().RemoveId(f.Id); err != nil {
		return err
	}

	path := f.FilePath()
	if err := os.Remove(path); err != nil {
		log.Printf("Failed to delete file: %s", path)
	}

	return nil
}

// Fetch from database: site that file belongs to
func (f *File) FindSite() *Site {
	return f.dbSession.FindSite(f.SiteId)
}

// Returns file path
func (f *File) FilePath() string {
	return core.UploadSiteFilePath(f.SiteId, f.Path)
}

// Returns file URL
func (f *File) URL() string {
	return core.UploadSiteUrlPath(f.SiteId, f.Path)
}
