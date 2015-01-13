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

	Logo  bson.ObjectId `bson:"logo,omitempty"  json:"logo,omitempty"`
	Cover bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`

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

type SitesList []*Site

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
	// @todo Remove that ?
	links := map[string]interface{}{
		"posts":   fmt.Sprintf("/api/sites/%s/posts", site.Id),
		"events":  fmt.Sprintf("/api/sites/%s/events", site.Id),
		"pages":   fmt.Sprintf("/api/sites/%s/pages", site.Id),
		"actions": fmt.Sprintf("/api/sites/%s/actions", site.Id),
		"images":  fmt.Sprintf("/api/sites/%s/images", site.Id),
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

// Returns the total number of posts
func (site *Site) PostsNb() int {
	result, err := site.postsBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// Fetch from database: all posts belonging to site
func (site *Site) FindPosts(skip int, limit int) *PostsList {
	result := PostsList{}

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

	// inject dbSession in all result items
	for _, post := range result {
		post.dbSession = site.dbSession
	}

	return &result
}

// Fetch from database: all events belonging to site
func (site *Site) FindEvents() *EventsList {
	result := EventsList{}

	if err := site.dbSession.EventsCol().Find(bson.M{"site_id": site.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

func (site *Site) pagesBaseQuery() *mgo.Query {
	return site.dbSession.PagesCol().Find(bson.M{"site_id": site.Id})
}

// Returns the total number of pages
func (site *Site) PagesNb() int {
	result, err := site.pagesBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// Fetch from database: all pages belonging to site
func (site *Site) FindPages(skip int, limit int) *PagesList {
	result := PagesList{}

	query := site.pagesBaseQuery().Sort("-created_at")

	if skip > 0 {
		query = query.Skip(skip)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.All(&result); err != nil {
		panic(err)
	}

	// inject dbSession in all result items
	for _, page := range result {
		page.dbSession = site.dbSession
	}

	return &result
}

// Fetch from database: all actions belonging to site
func (site *Site) FindActions() *ActionsList {
	result := ActionsList{}

	if err := site.dbSession.ActionsCol().Find(bson.M{"site_id": site.Id}).All(&result); err != nil {
		panic(err)
	}

	// @todo Inject dbSession in all result items

	return &result
}

func (site *Site) imagesBaseQuery() *mgo.Query {
	return site.dbSession.ImagesCol().Find(bson.M{"site_id": site.Id})
}

func (site *Site) ImagesNb() int {
	result, err := site.imagesBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// Fetch from database: all images belonging to site
func (site *Site) FindImages(skip int, limit int) *ImagesList {
	result := ImagesList{}

	query := site.imagesBaseQuery().Sort("-created_at")

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

// Fetch Logo from database
func (site *Site) FindLogo() *Image {
	if site.Logo != "" {
		var result Image

		if err := site.dbSession.ImagesCol().FindId(site.Logo).One(&result); err != nil {
			return nil
		}

		result.dbSession = site.dbSession

		return &result
	}

	return nil
}

// Fetch Cover from database
func (site *Site) FindCover() *Image {
	if site.Cover != "" {
		var result Image

		if err := site.dbSession.ImagesCol().FindId(site.Cover).One(&result); err != nil {
			return nil
		}

		result.dbSession = site.dbSession

		return &result
	}

	return nil
}

// Remove all references to given image from database
func (site *Site) RemoveImageReferences(image *Image) error {
	// remove image reference from site settings
	fieldsToDelete := []string{}

	if site.Logo == image.Id {
		fieldsToDelete = append(fieldsToDelete, "logo")
	}

	if site.Cover == image.Id {
		fieldsToDelete = append(fieldsToDelete, "cover")
	}

	if len(fieldsToDelete) > 0 {
		site.DeleteFields(fieldsToDelete)
	}

	// remove image references from posts
	if err := site.dbSession.RemoveImageReferencesFromPosts(image); err != nil {
		return err
	}

	// remove image references from pages
	if err := site.dbSession.RemoveImageReferencesFromPages(image); err != nil {
		return err
	}

	// @todo remove image references from actions / events / members / pages

	return nil
}

// Update site in database
func (site *Site) Update(newSite *Site) error {
	var set, unset, modifier bson.D

	if site.Name != newSite.Name {
		site.Name = newSite.Name

		if site.Name == "" {
			unset = append(unset, bson.DocElem{"name", 1})
		} else {
			set = append(set, bson.DocElem{"name", site.Name})
		}
	}

	if site.Tagline != newSite.Tagline {
		site.Tagline = newSite.Tagline

		if site.Tagline == "" {
			unset = append(unset, bson.DocElem{"tagline", 1})
		} else {
			set = append(set, bson.DocElem{"tagline", site.Tagline})
		}
	}

	if site.Description != newSite.Description {
		site.Description = newSite.Description

		if site.Description == "" {
			unset = append(unset, bson.DocElem{"description", 1})
		} else {
			set = append(set, bson.DocElem{"description", site.Description})
		}
	}

	if site.MoreDesc != newSite.MoreDesc {
		site.MoreDesc = newSite.MoreDesc

		if site.MoreDesc == "" {
			unset = append(unset, bson.DocElem{"more_desc", 1})
		} else {
			set = append(set, bson.DocElem{"more_desc", site.MoreDesc})
		}
	}

	if site.JoinText != newSite.JoinText {
		site.JoinText = newSite.JoinText

		if site.JoinText == "" {
			unset = append(unset, bson.DocElem{"join_text", 1})
		} else {
			set = append(set, bson.DocElem{"join_text", site.JoinText})
		}
	}

	if site.Logo != newSite.Logo {
		site.Logo = newSite.Logo

		if site.Logo == "" {
			unset = append(unset, bson.DocElem{"logo", 1})
		} else {
			set = append(set, bson.DocElem{"logo", site.Logo})
		}
	}

	if site.Cover != newSite.Cover {
		site.Cover = newSite.Cover

		if site.Cover == "" {
			unset = append(unset, bson.DocElem{"cover", 1})
		} else {
			set = append(set, bson.DocElem{"cover", site.Cover})
		}
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		return site.dbSession.SitesCol().UpdateId(site.Id, modifier)
	} else {
		return nil
	}
}

// Delete site fields from database
func (site *Site) DeleteFields(fields []string) error {
	var unset bson.D

	for _, field := range fields {
		unset = append(unset, bson.DocElem{field, 1})
	}

	return site.dbSession.SitesCol().UpdateId(site.Id, bson.M{"$unset": unset})
}
