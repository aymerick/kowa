package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SITES_COL_NAME = "sites"
)

type SitePageSettings struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Kind    string        `bson:"kind"    json:"kind"` // 'contact' || 'actions' || 'posts' || 'events' || 'staff'
	Title   string        `bson:"title"   json:"title"`
	Tagline string        `bson:"tagline" json:"tagline"`
	// @todo Photo
}

type Site struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UserId    bson.ObjectId `bson:"user_id"       json:"user"`

	Name        string `bson:"name"        json:"name"`
	Tagline     string `bson:"tagline"     json:"tagline"`
	Description string `bson:"description" json:"description"`
	MoreDesc    string `bson:"more_desc"   json:"moreDesc"`
	JoinText    string `bson:"join_text"   json:"joinText"`
	// @todo Logo
	// @todo Photo

	PageSettings []SitePageSettings `bson:"page_settings" json:"pageSettings"`

	// @todo Address
	// @todo Email
	// @todo Facebook
	// @todo Twitter
	// @todo GooglePlus
}

type SitesList []Site

// Sites collection
func SitesCol() *mgo.Collection {
	return DB().C(SITES_COL_NAME)
}

func AllSites() *SitesList {
	var result SitesList

	SitesCol().Find(nil).All(&result)

	return &result
}
