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
	// FileMembership represents a membership file kind
	FileMembership = "membership"

	filesColName = "files"
)

// FileKinds all possible file kinds
var FileKinds = []string{FileMembership}

// File represents a file
type File struct {
	dbSession *DBSession `bson:"-"`

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

// FilesList represents a list of files
type FilesList []*File

// FileJSON represents the json version of a file
type FileJSON struct {
	File
	URL string `json:"url"`
}

// IsValidFileKind returns true if argument is a file kind
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

// FilesCol returns files collection
func (session *DBSession) FilesCol() *mgo.Collection {
	return session.DB().C(filesColName)
}

// EnsureFilesIndexes ensure indexes on files collection
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

// FindFile finds a file by id
func (session *DBSession) FindFile(fileID bson.ObjectId) *File {
	var result File

	if err := session.FilesCol().FindId(fileID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreateFile creates a new file in database
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

// MarshalJSON implements the json.Marshaler interface
func (f *File) MarshalJSON() ([]byte, error) {

	fileJSON := FileJSON{
		File: *f,
		URL:  f.URL(),
	}

	return json.Marshal(fileJSON)
}

// Delete deletes file from database
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

// FindSite fetches site that file belongs to
func (f *File) FindSite() *Site {
	return f.dbSession.FindSite(f.SiteId)
}

// FilePath returns file path
func (f *File) FilePath() string {
	return core.UploadSiteFilePath(f.SiteId, f.Path)
}

// URL returns file URL
func (f *File) URL() string {
	return core.UploadSiteUrlPath(f.SiteId, f.Path)
}
