package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SITES_COL_NAME = "sites"
)

type SitePageSettings struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Kind    string        `bson:"kind"    json:"kind"` // 'contact' || 'actions' || 'posts' || 'events' || 'staff'
	Title   string        `bson:"title"   json:"title"`
	Tagline string        `bson:"tagline" json:"tagline"`
	// @todo Photo
}

type Site struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UserId    string        `bson:"user_id"       json:"user"`

	Name        string `bson:"name"        json:"name"`
	Tagline     string `bson:"tagline"     json:"tagline"`
	Description string `bson:"description" json:"description"`
	MoreDesc    string `bson:"more_desc"   json:"moreDesc"`
	JoinText    string `bson:"join_text"   json:"joinText"`
	// @todo Logo
	// @todo Photo

	PageSettings []SitePageSettings `bson:"page_settings" json:"pageSettings"`

	// @todo Address
	// @todo Email
	// @todo Facebook
	// @todo Twitter
	// @todo GooglePlus
}

type SiteJson struct {
	Site
	Links map[string]interface{} `json:"links"`
}

type SitesList []Site

//
// DBSession
//

// Sites collection
func (session *DBSession) SitesCol() *mgo.Collection {
	return session.DB().C(SITES_COL_NAME)
}

// Ensure indexes on Users collection
func (session *DBSession) EnsureSitesIndexes() {
	index := mgo.Index{
		Key:        []string{"user_id"},
		Background: true,
	}

	err := session.SitesCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// Find site by id
func (session *DBSession) FindSite(siteId bson.ObjectId) *Site {
	var result Site

	if err := session.SitesCol().FindId(siteId).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

//
// Site
//

// Implements json.MarshalJSON
func (this *Site) MarshalJSON() ([]byte, error) {
	// inject 'links' needed by Ember Data
	links := map[string]interface{}{
		"posts":   fmt.Sprintf("/api/sites/%s/posts", this.Id.Hex()),
		"events":  fmt.Sprintf("/api/sites/%s/events", this.Id.Hex()),
		"pages":   fmt.Sprintf("/api/sites/%s/pages", this.Id.Hex()),
		"actions": fmt.Sprintf("/api/sites/%s/actions", this.Id.Hex()),
	}

	siteJson := SiteJson{
		Site:  *this,
		Links: links,
	}

	return json.Marshal(siteJson)
}

// Fetch from database: all posts belonging to site
func (this *Site) FindPosts(skip int, limit int) *PostsList {
	var result PostsList

	query := this.dbSession.PostsCol().Find(bson.M{"site_id": this.Id}).Sort("created_at")

	if skip > 0 {
		query = query.Skip(skip)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

// Fetch from database: all events belonging to site
func (this *Site) FindEvents() *EventsList {
	var result EventsList

	if err := this.dbSession.EventsCol().Find(bson.M{"site_id": this.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

// Fetch from database: all pages belonging to site
func (this *Site) FindPages() *PagesList {
	var result PagesList

	if err := this.dbSession.PagesCol().Find(bson.M{"site_id": this.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

// Fetch from database: all actions belonging to site
func (this *Site) FindActions() *ActionsList {
	var result ActionsList

	if err := this.dbSession.ActionsCol().Find(bson.M{"site_id": this.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}
