package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MEMBERS_COL_NAME = "members"
)

type Member struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site"`

	Fullname    string `bson:"full_name"   json:"fullName"`
	Role        string `bson:"role"        json:"role"`
	Description string `bson:"description" json:"description"`
	// @todo Photo
}

type MembersList []*Member

//
// DBSession
//

// Member collection
func (session *DBSession) MembersCol() *mgo.Collection {
	return session.DB().C(MEMBERS_COL_NAME)
}
