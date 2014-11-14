package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ACTIONS_COL_NAME = "actions"
)

type Action struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updated_at"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site_id"`

	Title string `bson:"title" json:"title"`
	Body  string `bson:"body"  json:"body"`
	// @todo Photos List
}

type ActionsList []Action

// Actions collection
func ActionsCol() *mgo.Collection {
	return DB().C(ACTIONS_COL_NAME)
}
