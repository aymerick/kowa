package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ACTIVITIES_COL_NAME = "activities"
)

type Activity struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	Title   string `bson:"title"   json:"title"`
	Summary string `bson:"summary" json:"summary"`
	Body    string `bson:"body"    json:"body"`

	Cover bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`

	// @todo Format (markdown|html)
}

type ActivitiesList []*Activity

//
// DBSession
//

// Activities collection
func (session *DBSession) ActivitiesCol() *mgo.Collection {
	return session.DB().C(ACTIVITIES_COL_NAME)
}

// Ensure indexes on Activities collection
func (session *DBSession) EnsureActivitiesIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.ActivitiesCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// Find activity by id
func (session *DBSession) FindActivity(activityId bson.ObjectId) *Activity {
	var result Activity

	if err := session.ActivitiesCol().FindId(activityId).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// Persists a new activity in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on activity record
func (session *DBSession) CreateActivity(activity *Activity) error {
	activity.Id = bson.NewObjectId()

	now := time.Now()
	activity.CreatedAt = now
	activity.UpdatedAt = now

	if err := session.ActivitiesCol().Insert(activity); err != nil {
		return err
	}

	activity.dbSession = session

	return nil
}

// Remove all references to given image from all activities
func (session *DBSession) RemoveImageReferencesFromActivities(image *Image) error {
	// @todo
	return nil
}

//
// Activity
//

// Fetch from database: site that activity belongs to
func (activity *Activity) FindSite() *Site {
	return activity.dbSession.FindSite(activity.SiteId)
}

// Fetch Cover from database
func (activity *Activity) FindCover() *Image {
	if activity.Cover != "" {
		var result Image

		if err := activity.dbSession.ImagesCol().FindId(activity.Cover).One(&result); err != nil {
			return nil
		}

		result.dbSession = activity.dbSession

		return &result
	}

	return nil
}

// Delete activity from database
func (activity *Activity) Delete() error {
	var err error

	// delete from database
	if err = activity.dbSession.ActivitiesCol().RemoveId(activity.Id); err != nil {
		return err
	}

	return nil
}

// Update activity in database
func (activity *Activity) Update(newActivity *Activity) (bool, error) {
	var set, unset, modifier bson.D

	if activity.Title != newActivity.Title {
		activity.Title = newActivity.Title

		if activity.Title == "" {
			unset = append(unset, bson.DocElem{"title", 1})
		} else {
			set = append(set, bson.DocElem{"title", activity.Title})
		}
	}

	if activity.Summary != newActivity.Summary {
		activity.Summary = newActivity.Summary

		if activity.Summary == "" {
			unset = append(unset, bson.DocElem{"summary", 1})
		} else {
			set = append(set, bson.DocElem{"summary", activity.Summary})
		}
	}

	if activity.Body != newActivity.Body {
		activity.Body = newActivity.Body

		if activity.Body == "" {
			unset = append(unset, bson.DocElem{"body", 1})
		} else {
			set = append(set, bson.DocElem{"body", activity.Body})
		}
	}

	if activity.Cover != newActivity.Cover {
		activity.Cover = newActivity.Cover

		if activity.Cover == "" {
			unset = append(unset, bson.DocElem{"cover", 1})
		} else {
			set = append(set, bson.DocElem{"cover", activity.Cover})
		}
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		activity.UpdatedAt = time.Now()
		set = append(set, bson.DocElem{"updated_at", activity.UpdatedAt})

		return true, activity.dbSession.ActivitiesCol().UpdateId(activity.Id, modifier)
	} else {
		return false, nil
	}
}
