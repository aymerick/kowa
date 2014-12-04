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
	Id      string `bson:"_id,omitempty" json:"id"`
	Kind    string `bson:"kind"    json:"kind"` // 'contact' || 'actions' || 'posts' || 'events' || 'staff'
	Title   string `bson:"title"   json:"title"`
	Tagline string `bson:"tagline" json:"tagline"`
	// @todo Photo
}

type Site struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        string    `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time `bson:"created_at"    json:"createdAt"`
	UserId    string    `bson:"user_id"       json:"user"`

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
func (session *DBSession) FindSite(siteId string) *Site {
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
func (site *Site) MarshalJSON() ([]byte, error) {
	// inject 'links' needed by Ember Data
	links := map[string]interface{}{
		"posts":   fmt.Sprintf("/api/sites/%s/posts", site.Id),
		"events":  fmt.Sprintf("/api/sites/%s/events", site.Id),
		"pages":   fmt.Sprintf("/api/sites/%s/pages", site.Id),
		"actions": fmt.Sprintf("/api/sites/%s/actions", site.Id),
	}

	siteJson := SiteJson{
		Site:  *site,
		Links: links,
	}

	return json.Marshal(siteJson)
}

func (site *Site) postsBaseQuery() *mgo.Query {
	return site.dbSession.PostsCol().Find(bson.M{"site_id": site.Id})
}

func (site *Site) PostsNb() int {
	result, err := site.postsBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// Fetch from database: all posts belonging to site
func (site *Site) FindPosts(skip int, limit int) *PostsList {
	var result PostsList

	query := site.postsBaseQuery().Sort("-created_at")

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
func (site *Site) FindEvents() *EventsList {
	var result EventsList

	if err := site.dbSession.EventsCol().Find(bson.M{"site_id": site.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

// Fetch from database: all pages belonging to site
func (site *Site) FindPages() *PagesList {
	var result PagesList

	if err := site.dbSession.PagesCol().Find(bson.M{"site_id": site.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

// Fetch from database: all actions belonging to site
func (site *Site) FindActions() *ActionsList {
	var result ActionsList

	if err := site.dbSession.ActionsCol().Find(bson.M{"site_id": site.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

// Update site in database
func (site *Site) Update(newSite *Site) error {
	fields := bson.M{}

	if site.Name != newSite.Name {
		site.Name = newSite.Name
		fields["name"] = site.Name
	}

	if site.Tagline != newSite.Tagline {
		site.Tagline = newSite.Tagline
		fields["tagline"] = site.Tagline
	}

	if site.Description != newSite.Description {
		site.Description = newSite.Description
		fields["description"] = site.Description
	}

	if site.MoreDesc != newSite.MoreDesc {
		site.MoreDesc = newSite.MoreDesc
		fields["more_desc"] = site.MoreDesc
	}

	if site.JoinText != newSite.JoinText {
		site.JoinText = newSite.JoinText
		fields["join_text"] = site.JoinText
	}

	if len(fields) > 0 {
		return site.dbSession.SitesCol().UpdateId(site.Id, bson.M{"$set": fields})
	} else {
		return nil
	}
}
