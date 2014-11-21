package models

import (
	"encoding/json"
	"fmt"
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

type UserJson struct {
	User
	Links map[string]interface{} `json:"links"` // needed by Ember Data
}

type UsersList []User

// Handler for Users collection
func UsersCol() *mgo.Collection {
	return DB().C(USERS_COL_NAME)
}

// Find user by id
func FindUser(userId string) *User {
	var result User

	// @todo Handle err
	UsersCol().FindId(bson.ObjectIdHex(userId)).One(&result)

	return &result
}

// Implements json.MarshalJSON
func (this *User) MarshalJSON() ([]byte, error) {
	// inject 'links' needed by Ember Data
	links := map[string]interface{}{"sites": fmt.Sprintf("/api/users/%s/sites", this.Id.Hex())}

	userJson := UserJson{
		User:  *this,
		Links: links,
	}

	return json.Marshal(userJson)
}

// Find all sites belonging to user
func (this *User) FindSites() *SitesList {
	var result SitesList

	SitesCol().Find(bson.M{"user_id": this.Id}).All(&result)

	return &result
}
