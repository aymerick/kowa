package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	EVENTS_COL_NAME = "events"
)

type Event struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	Date  time.Time `bson:"date"  json:"date"`
	Title string    `bson:"title" json:"title"`
	Body  string    `bson:"body"  json:"body"`
	Place string    `bson:"place" json:"place"`
	// @todo Photo
}

type EventsList []*Event

//
// DBSession
//

// Events collection
func (session *DBSession) EventsCol() *mgo.Collection {
	return session.DB().C(EVENTS_COL_NAME)
}

// Ensure indexes on Events collection
func (session *DBSession) EnsureEventsIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.EventsCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

//
// Event
//

// @todo
