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

	// signed up but not confirmed
	USER_STATUS_PENDING = "pending"

	// signed up and confirmed
	USER_STATUS_ACTIVE = "active"
)

type User struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id          string    `bson:"_id,omitempty" json:"id"`
	CreatedAt   time.Time `bson:"created_at"    json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at"    json:"updatedAt"`
	ValidatedAt time.Time `bson:"validated_at"  json:"validatedAt"`
	Admin       bool      `bson:"admin"         json:"admin"`
	Status      string    `bson:"status"        json:"status"`

	Email     string `bson:"email"      json:"email"`
	FirstName string `bson:"first_name" json:"firstName"`
	LastName  string `bson:"last_name"  json:"lastName"`
	Lang      string `bson:"lang"       json:"lang"`
	Password  string `bson:"password"   json:"-"`
}

type UserJson struct {
	User
	Links map[string]interface{} `json:"links"`
}

type UsersList []*User

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

// Persists a new user in database
// Side effect: 'CreatedAt' and 'UpdatedAt' fields are set on user record
func (session *DBSession) CreateUser(user *User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := session.UsersCol().Insert(user); err != nil {
		return err
	}

	user.dbSession = session

	return nil
}

//
// User
//

// Returns true if user have an active account
func (user *User) Active() bool {
	return user.Status == USER_STATUS_ACTIVE
}

// Returns true if user account has been validated
func (user *User) AccountValidated() bool {
	return !user.ValidatedAt.IsZero()
}

// Returns user fullname
func (user *User) FullName() string {
	if user.FirstName == "" {
		return user.LastName
	} else if user.LastName == "" {
		return user.FirstName
	} else {
		return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
}

// Returns user display name, usefull if fullname is empty
func (user *User) DisplayName() string {
	result := user.FullName()

	if result == "" {
		result = user.Id
	}

	return result
}

// Returns user mail address with format: User Name <email@addre.ss>
func (user *User) MailAddress() string {
	return fmt.Sprintf("%s <%s>", user.DisplayName(), user.Email)
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
	result := SitesList{}

	// @todo Handle err
	user.dbSession.SitesCol().Find(bson.M{"user_id": user.Id}).All(&result)

	for _, site := range result {
		site.dbSession = user.dbSession
	}

	return &result
}

// Update user in database
func (user *User) Update(newUser *User) (bool, error) {
	var set, unset, modifier bson.D

	// FirstName
	if user.FirstName != newUser.FirstName {
		user.FirstName = newUser.FirstName

		if user.FirstName == "" {
			unset = append(unset, bson.DocElem{"first_name", 1})
		} else {
			set = append(set, bson.DocElem{"first_name", user.FirstName})
		}
	}

	// LastName
	if user.LastName != newUser.LastName {
		user.LastName = newUser.LastName

		if user.LastName == "" {
			unset = append(unset, bson.DocElem{"last_name", 1})
		} else {
			set = append(set, bson.DocElem{"last_name", user.LastName})
		}
	}

	// Lang
	if user.Lang != newUser.Lang {
		user.Lang = newUser.Lang

		if user.Lang == "" {
			unset = append(unset, bson.DocElem{"lang", 1})
		} else {
			set = append(set, bson.DocElem{"lang", user.Lang})
		}
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		user.UpdatedAt = time.Now()
		set = append(set, bson.DocElem{"updated_at", user.UpdatedAt})

		return true, user.dbSession.UsersCol().UpdateId(user.Id, modifier)
	} else {
		return false, nil
	}
}

func (user *User) SetValues(values bson.M) error {
	// @todo Set UpdatedAt field
	return user.dbSession.UsersCol().UpdateId(user.Id, bson.D{{"$set", values}})
}

// Set user account as validated
func (user *User) SetAccountValidated() error {
	now := time.Now()

	fields := bson.M{
		"validated_at": now,
		"status":       USER_STATUS_ACTIVE,
	}

	if err := user.SetValues(fields); err != nil {
		return err
	}

	user.ValidatedAt = now
	user.Status = USER_STATUS_ACTIVE

	return nil
}
