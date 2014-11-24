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
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site"`

	Fullname    string `bson:"full_name"   json:"fullName"`
	Role        string `bson:"role"        json:"role"`
	Description string `bson:"description" json:"description"`
	// @todo Photo
}

type StaffMembersList []StaffMember

//
// DBSession
//

// StaffMember collection
func (session *DBSession) StaffMembersCol() *mgo.Collection {
	return session.DB().C(STAFF_MEMBERS_COL_NAME)
}
