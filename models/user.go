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
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`

	FirstName string `bson:"first_name" json:"firstName"`
	LastName  string `bson:"last_name"  json:"lastName"`
}

type UsersList []User

// Users collection
func UsersCol() *mgo.Collection {
	return DB().C(USERS_COL_NAME)
}

func FindUser(userId string) *User {
	var result User

	// @todo Handle err
	UsersCol().FindId(bson.ObjectIdHex(userId)).One(&result)

	return &result
}
