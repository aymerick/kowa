package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	eventsColName = "events"
)

// Event represents an event
type Event struct {
	dbSession *DBSession `bson:"-"`

	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteID    string        `bson:"site_id"       json:"site"`

	StartDate time.Time     `bson:"start_date"      json:"startDate,omitempty"`
	EndDate   time.Time     `bson:"end_date"        json:"endDate,omitempty"`
	Title     string        `bson:"title"           json:"title"`
	Body      string        `bson:"body"            json:"body"`
	Format    string        `bson:"format"          json:"format"`
	Place     string        `bson:"place"           json:"place"`
	Cover     bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`
}

// EventsList represents a list of events
type EventsList []*Event

//
// DBSession
//

// EventsCol returns the events collection
func (session *DBSession) EventsCol() *mgo.Collection {
	return session.DB().C(eventsColName)
}

// EnsureEventsIndexes ensures indexes on events collection
func (session *DBSession) EnsureEventsIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.EventsCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// FindEvent finds an event by id
func (session *DBSession) FindEvent(eventID bson.ObjectId) *Event {
	var result Event

	if err := session.EventsCol().FindId(eventID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreateEvent creates a new event in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on event record
func (session *DBSession) CreateEvent(event *Event) error {
	event.ID = bson.NewObjectId()

	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	if err := session.EventsCol().Insert(event); err != nil {
		return err
	}

	event.dbSession = session

	return nil
}

// RemoveImageReferencesFromEvents remove all references to given image from all events
func (session *DBSession) RemoveImageReferencesFromEvents(image *Image) error {
	// @todo
	return nil
}

//
// Event
//

// FindSite fetch site that event belongs to
func (event *Event) FindSite() *Site {
	return event.dbSession.FindSite(event.SiteID)
}

// FindCover fetches cover from database
func (event *Event) FindCover() *Image {
	if event.Cover != "" {
		var result Image

		if err := event.dbSession.ImagesCol().FindId(event.Cover).One(&result); err != nil {
			return nil
		}

		result.dbSession = event.dbSession

		return &result
	}

	return nil
}

// Delete deletes event from database
func (event *Event) Delete() error {
	var err error

	// delete from database
	if err = event.dbSession.EventsCol().RemoveId(event.ID); err != nil {
		return err
	}

	return nil
}

// Update updates event in database
func (event *Event) Update(newEvent *Event) (bool, error) {
	var set, unset, modifier bson.D

	// Startdate
	if event.StartDate != newEvent.StartDate {
		event.StartDate = newEvent.StartDate

		if event.StartDate.IsZero() {
			unset = append(unset, bson.DocElem{"start_date", 1})
		} else {
			set = append(set, bson.DocElem{"start_date", event.StartDate})
		}
	}

	// Enddate
	if event.EndDate != newEvent.EndDate {
		event.EndDate = newEvent.EndDate

		if event.EndDate.IsZero() {
			unset = append(unset, bson.DocElem{"end_date", 1})
		} else {
			set = append(set, bson.DocElem{"end_date", event.EndDate})
		}
	}

	// Title
	if event.Title != newEvent.Title {
		event.Title = newEvent.Title

		if event.Title == "" {
			unset = append(unset, bson.DocElem{"title", 1})
		} else {
			set = append(set, bson.DocElem{"title", event.Title})
		}
	}

	// Body
	if event.Body != newEvent.Body {
		event.Body = newEvent.Body

		if event.Body == "" {
			unset = append(unset, bson.DocElem{"body", 1})
		} else {
			set = append(set, bson.DocElem{"body", event.Body})
		}
	}

	// Format
	newFormat := newEvent.Format
	if newFormat == "" {
		newFormat = DefaultFormat
	}

	if event.Format != newFormat {
		event.Format = newFormat

		set = append(set, bson.DocElem{"format", event.Format})
	}

	// Place
	if event.Place != newEvent.Place {
		event.Place = newEvent.Place

		if event.Place == "" {
			unset = append(unset, bson.DocElem{"place", 1})
		} else {
			set = append(set, bson.DocElem{"place", event.Place})
		}
	}

	// Cover
	if event.Cover != newEvent.Cover {
		event.Cover = newEvent.Cover

		if event.Cover == "" {
			unset = append(unset, bson.DocElem{"cover", 1})
		} else {
			set = append(set, bson.DocElem{"cover", event.Cover})
		}
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		event.UpdatedAt = time.Now()
		set = append(set, bson.DocElem{"updated_at", event.UpdatedAt})

		return true, event.dbSession.EventsCol().UpdateId(event.ID, modifier)
	}

	return false, nil
}
