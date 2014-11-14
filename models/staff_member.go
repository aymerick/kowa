package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	STAFF_MEMBERS_COL_NAME = "staff_members"
)

type StaffMember struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updated_at"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site_id"`

	Fullname    string `bson:"full_name"   json:"full_name"`
	Role        string `bson:"role"        json:"role"`
	Description string `bson:"description" json:"description"`
	// @todo Photo
}

type StaffMembersList []StaffMember

// StaffMember collection
func StaffMembersCol() *mgo.Collection {
	return DB().C(STAFF_MEMBERS_COL_NAME)
}
