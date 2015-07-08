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
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	Fullname    string        `bson:"fullname"        json:"fullname"`
	Role        string        `bson:"role"            json:"role"`
	Description string        `bson:"description"     json:"description"`
	Photo       bson.ObjectId `bson:"photo,omitempty" json:"photo,omitempty"`
	Order       int           `bson:"order"           json:"order"`
}

type MembersList []*Member

//
// DBSession
//

// Member collection
func (session *DBSession) MembersCol() *mgo.Collection {
	return session.DB().C(MEMBERS_COL_NAME)
}

// Ensure indexes on Members collection
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

// Find member by id
func (session *DBSession) FindMember(memberId bson.ObjectId) *Member {
	var result Member

	if err := session.MembersCol().FindId(memberId).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// Persists a new member in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on member record
func (session *DBSession) CreateMember(member *Member) error {
	member.Id = bson.NewObjectId()

	now := time.Now()
	member.CreatedAt = now
	member.UpdatedAt = now

	if err := session.MembersCol().Insert(member); err != nil {
		return err
	}

	member.dbSession = session

	return nil
}

// Remove all references to given image from all members
func (session *DBSession) RemoveImageReferencesFromMembers(image *Image) error {
	// @todo
	return nil
}

//
// Member
//

// Fetch from database: site that member belongs to
func (member *Member) FindSite() *Site {
	return member.dbSession.FindSite(member.SiteId)
}

// Fetch Photo from database
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

// Delete member from database
func (member *Member) Delete() error {
	var err error

	// delete from database
	if err = member.dbSession.MembersCol().RemoveId(member.Id); err != nil {
		return err
	}

	return nil
}

// Update member in database
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

		return true, member.dbSession.MembersCol().UpdateId(member.Id, modifier)
	} else {
		return false, nil
	}
}
