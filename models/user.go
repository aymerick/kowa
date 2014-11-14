package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	USERS_COL_NAME = "users"
)

type User struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`

	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name"  json:"last_name"`
}

type UsersList []User

// Users collection
func UsersCol() *mgo.Collection {
	return DB().C(USERS_COL_NAME)
}

func AllUsers() *UsersList {
	var result UsersList

	UsersCol().Find(nil).All(&result)

	return &result
}
