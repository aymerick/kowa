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
	dbSession *DBSession `bson:"-" json:"-"`

	Id        string    `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time `bson:"created_at"    json:"createdAt"`

	Email     string `bson:"email"      json:"email"`
	FirstName string `bson:"first_name" json:"firstName"`
	LastName  string `bson:"last_name"  json:"lastName"`
	Password  string `bson:"password"   json:"-"`
}

type UserJson struct {
	User
	Links map[string]interface{} `json:"links"`
}

type UsersList []User

//
// DBSession
//

// Handler for Users collection
func (session *DBSession) UsersCol() *mgo.Collection {
	return session.DB().C(USERS_COL_NAME)
}

// Ensure indexes on Users collection
func (session *DBSession) EnsureUsersIndexes() {
	// Find by email
	index := mgo.Index{
		Key:        []string{"email"},
		Background: true,
	}

	err := session.UsersCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// Find user by id
func (session *DBSession) FindUser(userId string) *User {
	var result User

	if err := session.UsersCol().FindId(userId).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// Find user by email
func (session *DBSession) FindUserByEmail(email string) *User {
	var result User

	if err := session.UsersCol().Find(bson.M{"email": email}).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

//
// User
//

func (user *User) Fullname() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

// Implements json.MarshalJSON
func (user *User) MarshalJSON() ([]byte, error) {
	// inject 'links' needed by Ember Data
	links := map[string]interface{}{"sites": fmt.Sprintf("/api/users/%s/sites", user.Id)}

	userJson := UserJson{
		User:  *user,
		Links: links,
	}

	return json.Marshal(userJson)
}

// Fetch from database: all sites belonging to user
func (user *User) FindSites() *SitesList {
	var result SitesList

	// @todo Handle err
	user.dbSession.SitesCol().Find(bson.M{"user_id": user.Id}).All(&result)

	// @todo Inject dbSession in all result sites

	return &result
}
