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
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updated_at"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site_id"`

	Date  time.Time `bson:"date"  json:"date"`
	Title string    `bson:"title" json:"title"`
	Body  string    `bson:"body"  json:"body"`
	Place string    `bson:"place" json:"place"`
	// @todo Photo
}

type EventsList []Event

// Events collection
func EventsCol() *mgo.Collection {
	return DB().C(EVENTS_COL_NAME)
}
