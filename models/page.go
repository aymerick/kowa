package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	PAGES_COL_NAME = "pages"
)

type Page struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site"`

	Title   string `bson:"title"   json:"title"`
	Tagline string `bson:"tagline" json:"tagline"`
	Body    string `bson:"body"    json:"body"`
	// @todo Photo
}

type PagesList []Page

//
// DBSession
//

// Pages collection
func (session *DBSession) PagesCol() *mgo.Collection {
	return session.DB().C(PAGES_COL_NAME)
}

// Ensure indexes on Pages collection
func (session *DBSession) EnsurePagesIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.PagesCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

//
// Page
//

// @todo
