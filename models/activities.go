package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ACTIVITIES_COL_NAME = "activities"
)

type Activity struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	Title string `bson:"title" json:"title"`
	Body  string `bson:"body"  json:"body"`
	// @todo Photos List
}

type ActivitiesList []*Activity

//
// DBSession
//

// Activities collection
func (session *DBSession) ActivitiesCol() *mgo.Collection {
	return session.DB().C(ACTIVITIES_COL_NAME)
}

// Ensure indexes on Activities collection
func (session *DBSession) EnsureActivitiesIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.ActivitiesCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

//
// Activity
//

// @todo
