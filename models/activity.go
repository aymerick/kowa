package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	activitesColName = "activities"
)

// Activity represents an activity
type Activity struct {
	dbSession *DBSession `bson:"-"`

	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteID    string        `bson:"site_id"       json:"site"`

	Title   string        `bson:"title"           json:"title"`
	Summary string        `bson:"summary"         json:"summary"`
	Body    string        `bson:"body"            json:"body"`
	Cover   bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`
}

// ActivitiesList holds a list of Activity
type ActivitiesList []*Activity

//
// DBSession
//

// ActivitiesCol returns the activities collection
func (session *DBSession) ActivitiesCol() *mgo.Collection {
	return session.DB().C(activitesColName)
}

// EnsureActivitiesIndexes ensures indexes on activities collection
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

// FindActivity finds an activity by id
func (session *DBSession) FindActivity(activityID bson.ObjectId) *Activity {
	var result Activity

	if err := session.ActivitiesCol().FindId(activityID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreateActivity persists a new activity in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on activity record
func (session *DBSession) CreateActivity(activity *Activity) error {
	activity.ID = bson.NewObjectId()

	now := time.Now()
	activity.CreatedAt = now
	activity.UpdatedAt = now

	if err := session.ActivitiesCol().Insert(activity); err != nil {
		return err
	}

	activity.dbSession = session

	return nil
}

// RemoveImageReferencesFromActivities removes all references to given image from all activities
func (session *DBSession) RemoveImageReferencesFromActivities(image *Image) error {
	// @todo
	return nil
}

//
// Activity
//

// FindSite fetches the site that activity belongs to
func (activity *Activity) FindSite() *Site {
	return activity.dbSession.FindSite(activity.SiteID)
}

// FindCover fetches activity cover
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

// Delete deletes activity from database
func (activity *Activity) Delete() error {
	var err error

	// delete from database
	if err = activity.dbSession.ActivitiesCol().RemoveId(activity.ID); err != nil {
		return err
	}

	return nil
}

// Update updates activity in database
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

		return true, activity.dbSession.ActivitiesCol().UpdateId(activity.ID, modifier)
	}

	return false, nil
}
