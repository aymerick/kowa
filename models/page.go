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

// Pages collection
func PagesCol() *mgo.Collection {
	return DB().C(PAGES_COL_NAME)
}
