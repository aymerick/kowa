package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	membersColName = "members"
)

// Member represents a member
type Member struct {
	dbSession *DBSession `bson:"-"`

	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteID    string        `bson:"site_id"       json:"site"`

	Fullname    string        `bson:"fullname"        json:"fullname"`
	Role        string        `bson:"role"            json:"role"`
	Description string        `bson:"description"     json:"description"`
	Photo       bson.ObjectId `bson:"photo,omitempty" json:"photo,omitempty"`
	Order       int           `bson:"order"           json:"order"`
}

// MembersList represents a list of Member
type MembersList []*Member

//
// DBSession
//

// MembersCol returns members collection
func (session *DBSession) MembersCol() *mgo.Collection {
	return session.DB().C(membersColName)
}

// EnsureMembersIndexes ensure indexes on members collection
func (session *DBSession) EnsureMembersIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id", "order"},
		Background: true,
	}

	err := session.MembersCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// FindMember finds member by id
func (session *DBSession) FindMember(memberID bson.ObjectId) *Member {
	var result Member

	if err := session.MembersCol().FindId(memberID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreateMember creates a new member in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on member record
func (session *DBSession) CreateMember(member *Member) error {
	member.ID = bson.NewObjectId()

	now := time.Now()
	member.CreatedAt = now
	member.UpdatedAt = now

	if err := session.MembersCol().Insert(member); err != nil {
		return err
	}

	member.dbSession = session

	return nil
}

// RemoveImageReferencesFromMembers removes all references to given image from all members
func (session *DBSession) RemoveImageReferencesFromMembers(image *Image) error {
	// @todo
	return nil
}

//
// Member
//

// FindSite fetches site that member belongs to
func (member *Member) FindSite() *Site {
	return member.dbSession.FindSite(member.SiteID)
}

// FindPhoto fetches Photo from database
func (member *Member) FindPhoto() *Image {
	if member.Photo != "" {
		var result Image

		if err := member.dbSession.ImagesCol().FindId(member.Photo).One(&result); err != nil {
			return nil
		}

		result.dbSession = member.dbSession

		return &result
	}

	return nil
}

// Delete deletes member from database
func (member *Member) Delete() error {
	var err error

	// delete from database
	if err = member.dbSession.MembersCol().RemoveId(member.ID); err != nil {
		return err
	}

	return nil
}

// Update updates member in database
func (member *Member) Update(newMember *Member) (bool, error) {
	var set, unset, modifier bson.D

	// Fullname
	if member.Fullname != newMember.Fullname {
		member.Fullname = newMember.Fullname

		if member.Fullname == "" {
			unset = append(unset, bson.DocElem{"fullname", 1})
		} else {
			set = append(set, bson.DocElem{"fullname", member.Fullname})
		}
	}

	// Role
	if member.Role != newMember.Role {
		member.Role = newMember.Role

		if member.Role == "" {
			unset = append(unset, bson.DocElem{"role", 1})
		} else {
			set = append(set, bson.DocElem{"role", member.Role})
		}
	}

	// Description
	if member.Description != newMember.Description {
		member.Description = newMember.Description

		if member.Description == "" {
			unset = append(unset, bson.DocElem{"description", 1})
		} else {
			set = append(set, bson.DocElem{"description", member.Description})
		}
	}

	// Photo
	if member.Photo != newMember.Photo {
		member.Photo = newMember.Photo

		if member.Photo == "" {
			unset = append(unset, bson.DocElem{"photo", 1})
		} else {
			set = append(set, bson.DocElem{"photo", member.Photo})
		}
	}

	// Order
	if member.Order != newMember.Order {
		member.Order = newMember.Order

		set = append(set, bson.DocElem{"order", member.Order})
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		member.UpdatedAt = time.Now()
		set = append(set, bson.DocElem{"updated_at", member.UpdatedAt})

		return true, member.dbSession.MembersCol().UpdateId(member.ID, modifier)
	}

	return false, nil
}
