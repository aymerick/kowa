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
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site"`

	Title string `bson:"title" json:"title"`
	Body  string `bson:"body"  json:"body"`
	// @todo Photos List
}

type ActionsList []Action

//
// DBSession
//

// Actions collection
func (session *DBSession) ActionsCol() *mgo.Collection {
	return session.DB().C(ACTIONS_COL_NAME)
}

// Ensure indexes on Actions collection
func (session *DBSession) EnsureActionsIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.ActionsCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

//
// Action
//

// @todo
