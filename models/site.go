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
	Kind    string `bson:"kind"    json:"kind"` // 'contact' || 'activities' || 'posts' || 'events' || 'members'
	Title   string `bson:"title"   json:"title"`
	Tagline string `bson:"tagline" json:"tagline"`
	// @todo Photo
}

type Site struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        string    `bson:"_id,omitempty"         json:"id"`
	CreatedAt time.Time `bson:"created_at"            json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at"            json:"updatedAt"`
	ChangedAt time.Time `bson:"changed_at,omitempty"  json:"changedAt,omitempty"`
	BuiltAt   time.Time `bson:"built_at,omitempty"    json:"builtAt,omitempty"`
	UserId    string    `bson:"user_id"               json:"user"`

	Name        string `bson:"name"        json:"name"`
	Tagline     string `bson:"tagline"     json:"tagline"`
	Description string `bson:"description" json:"description"`
	MoreDesc    string `bson:"more_desc"   json:"moreDesc"`
	JoinText    string `bson:"join_text"   json:"joinText"`

	Email   string `bson:"email"   json:"email"`
	Address string `bson:"address" json:"address"`

	Facebook   string `bson:"facebook"    json:"facebook"`
	Twitter    string `bson:"twitter"     json:"twitter"`
	GooglePlus string `bson:"google_plus" json:"googlePlus"`

	Logo  bson.ObjectId `bson:"logo,omitempty"  json:"logo,omitempty"`
	Cover bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`

	PageSettings []SitePageSettings `bson:"page_settings" json:"pageSettings"`

	// build settings
	Theme   string `bson:"theme"    json:"theme"`
	UglyURL bool   `bson:"ugly_url" json:"uglyUrl"`
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

// Persists a new site in database
// Side effect: 'CreatedAt' and 'UpdatedAt' fields are set on site record
func (session *DBSession) CreateSite(site *Site) error {
	now := time.Now()
	site.CreatedAt = now
	site.UpdatedAt = now

	if err := session.SitesCol().Insert(site); err != nil {
		return err
	}

	site.dbSession = session

	return nil
}

//
// Site
//

// Implements json.MarshalJSON
func (site *Site) MarshalJSON() ([]byte, error) {
	// inject 'links' needed by Ember Data
	// @todo Remove that ?
	links := map[string]interface{}{
		"posts":      fmt.Sprintf("/api/sites/%s/posts", site.Id),
		"events":     fmt.Sprintf("/api/sites/%s/events", site.Id),
		"pages":      fmt.Sprintf("/api/sites/%s/pages", site.Id),
		"activities": fmt.Sprintf("/api/sites/%s/activities", site.Id),
		"members":    fmt.Sprintf("/api/sites/%s/members", site.Id),
		"images":     fmt.Sprintf("/api/sites/%s/images", site.Id),
	}

	siteJson := SiteJson{
		Site:  *site,
		Links: links,
	}

	return json.Marshal(siteJson)
}

//
// Site posts
//

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

func (site *Site) FindAllPosts() *PostsList {
	return site.FindPosts(0, 0)
}

//
// Site events
//

func (site *Site) eventsBaseQuery() *mgo.Query {
	return site.dbSession.EventsCol().Find(bson.M{"site_id": site.Id})
}

// Returns the total number of events
func (site *Site) EventsNb() int {
	result, err := site.eventsBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// Fetch from database: all events belonging to site
func (site *Site) FindEvents(skip int, limit int) *EventsList {
	result := EventsList{}

	query := site.eventsBaseQuery().Sort("-start_date")

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
	for _, event := range result {
		event.dbSession = site.dbSession
	}

	return &result
}

func (site *Site) FindAllEvents() *EventsList {
	return site.FindEvents(0, 0)
}

//
// Site pages
//

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

	query := site.pagesBaseQuery().Sort("created_at")

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

func (site *Site) FindAllPages() *PagesList {
	return site.FindPages(0, 0)
}

//
// Site activities
//

func (site *Site) activitiesBaseQuery() *mgo.Query {
	return site.dbSession.ActivitiesCol().Find(bson.M{"site_id": site.Id})
}

// Returns the total number of activities
func (site *Site) ActivitiesNb() int {
	result, err := site.activitiesBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// Fetch from database: all activities belonging to site
func (site *Site) FindActivities(skip int, limit int) *ActivitiesList {
	result := ActivitiesList{}

	query := site.activitiesBaseQuery().Sort("created_at")

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
	for _, activity := range result {
		activity.dbSession = site.dbSession
	}

	return &result
}

func (site *Site) FindAllActivities() *ActivitiesList {
	return site.FindActivities(0, 0)
}

//
// Site members
//

func (site *Site) membersBaseQuery() *mgo.Query {
	return site.dbSession.MembersCol().Find(bson.M{"site_id": site.Id})
}

// Returns the total number of members
func (site *Site) MembersNb() int {
	result, err := site.membersBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// Fetch from database: all members belonging to site
func (site *Site) FindMembers(skip int, limit int) *MembersList {
	result := MembersList{}

	query := site.membersBaseQuery().Sort("role")

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
	for _, member := range result {
		member.dbSession = site.dbSession
	}

	return &result
}

func (site *Site) FindAllMembers() *MembersList {
	return site.FindMembers(0, 0)
}

//
// Site images
//

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

func (site *Site) FindAllImages() *ImagesList {
	return site.FindImages(0, 0)
}

//
// Site fields
//

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

	// remove image references from events
	if err := site.dbSession.RemoveImageReferencesFromEvents(image); err != nil {
		return err
	}

	// remove image references from pages
	if err := site.dbSession.RemoveImageReferencesFromPages(image); err != nil {
		return err
	}

	// remove image references from activities
	if err := site.dbSession.RemoveImageReferencesFromActivities(image); err != nil {
		return err
	}

	// remove image references from activities
	if err := site.dbSession.RemoveImageReferencesFromMembers(image); err != nil {
		return err
	}

	return nil
}

// Update site in database
func (site *Site) Update(newSite *Site) (bool, error) {
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

	if site.Email != newSite.Email {
		site.Email = newSite.Email

		if site.Email == "" {
			unset = append(unset, bson.DocElem{"email", 1})
		} else {
			set = append(set, bson.DocElem{"email", site.Email})
		}
	}

	if site.Address != newSite.Address {
		site.Address = newSite.Address

		if site.Address == "" {
			unset = append(unset, bson.DocElem{"address", 1})
		} else {
			set = append(set, bson.DocElem{"address", site.Address})
		}
	}

	if site.Facebook != newSite.Facebook {
		site.Facebook = newSite.Facebook

		if site.Facebook == "" {
			unset = append(unset, bson.DocElem{"facebook", 1})
		} else {
			set = append(set, bson.DocElem{"facebook", site.Facebook})
		}
	}

	if site.Twitter != newSite.Twitter {
		site.Twitter = newSite.Twitter

		if site.Twitter == "" {
			unset = append(unset, bson.DocElem{"twitter", 1})
		} else {
			set = append(set, bson.DocElem{"twitter", site.Twitter})
		}
	}

	if site.GooglePlus != newSite.GooglePlus {
		site.GooglePlus = newSite.GooglePlus

		if site.GooglePlus == "" {
			unset = append(unset, bson.DocElem{"google_plus", 1})
		} else {
			set = append(set, bson.DocElem{"google_plus", site.GooglePlus})
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

	if site.Theme != newSite.Theme {
		site.Theme = newSite.Theme

		if site.Theme == "" {
			unset = append(unset, bson.DocElem{"theme", 1})
		} else {
			set = append(set, bson.DocElem{"theme", site.Theme})
		}
	}

	if site.UglyURL != newSite.UglyURL {
		site.UglyURL = newSite.UglyURL

		if site.UglyURL == false {
			unset = append(unset, bson.DocElem{"ugly_url", 1})
		} else {
			set = append(set, bson.DocElem{"ugly_url", site.UglyURL})
		}
	}

	if (len(unset) > 0) || (len(set) > 0) {
		site.UpdatedAt = time.Now()
		set = append(set, bson.DocElem{"updated_at", site.UpdatedAt})
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		return true, site.dbSession.SitesCol().UpdateId(site.Id, modifier)
	} else {
		return false, nil
	}
}

func (site *Site) SetValues(values bson.D) error {
	return site.dbSession.SitesCol().UpdateId(site.Id, bson.D{{"$set", values}})
}

// Set the ChangedAt value
func (site *Site) SetChangedAt(value time.Time) error {
	if err := site.SetValues(bson.D{{"changed_at", value}}); err != nil {
		return err
	}

	site.ChangedAt = value
	return nil
}

// Set the BuiltAt value
func (site *Site) SetBuiltAt(value time.Time) error {
	if err := site.SetValues(bson.D{{"built_at", value}}); err != nil {
		return err
	}

	site.BuiltAt = value
	return nil
}

// Delete site fields from database
func (site *Site) DeleteFields(fields []string) error {
	var unset bson.D

	for _, field := range fields {
		unset = append(unset, bson.DocElem{field, 1})
	}

	return site.dbSession.SitesCol().UpdateId(site.Id, bson.M{"$unset": unset})
}
