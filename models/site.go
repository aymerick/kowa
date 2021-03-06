package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/themes"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	sitesColName = "sites"
)

// SitePagesSettingsKinds holds all possible page settings kinds
var SitePagesSettingsKinds map[string]bool

// SitePageSettings represents all the settings for a built-in page
type SitePageSettings struct {
	ID       bson.ObjectId `bson:"_id,omitempty"   json:"id"`
	Kind     string        `bson:"kind"            json:"kind"` // cf. SitePagesSettingsKinds
	Title    string        `bson:"title"           json:"title"`
	Tagline  string        `bson:"tagline"         json:"tagline"`
	Cover    bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`
	Disabled bool          `bson:"disabled"        json:"disabled"`
}

// SiteThemeSettings represents settings for given theme
type SiteThemeSettings struct {
	ID      bson.ObjectId     `bson:"_id,omitempty" json:"id"`
	Palette string            `bson:"palette,omitempty" json:"palette"`
	Custom  map[string]string `bson:"custom,omitempty" json:"-"`
	Theme   string            `bson:"-" json:"theme"`
}

// SiteThemeSettingsJSON is the JSON representation of SiteThemeSettings
type SiteThemeSettingsJSON struct {
	SiteThemeSettings

	// overrides the Custom field to provide a JSON string instead of an array of embedded documents
	Custom string `json:"custom,omitempty"`
}

// Site represents a site
type Site struct {
	dbSession *DBSession `bson:"-"`

	ID        string    `bson:"_id,omitempty"         json:"id"`
	CreatedAt time.Time `bson:"created_at"            json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at"            json:"updatedAt"`
	ChangedAt time.Time `bson:"changed_at,omitempty"  json:"changedAt,omitempty"`
	BuiltAt   time.Time `bson:"built_at,omitempty"    json:"builtAt,omitempty"`
	UserID    string    `bson:"user_id"               json:"user"`
	Lang      string    `bson:"lang"                  json:"lang"`
	TZ        string    `bson:"tz"                    json:"tz"`

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

	GoogleAnalytics string `bson:"google_analytics" json:"googleAnalytics"`

	// images
	Logo    bson.ObjectId `bson:"logo,omitempty"    json:"logo,omitempty"`
	Cover   bson.ObjectId `bson:"cover,omitempty"   json:"cover,omitempty"`
	Favicon bson.ObjectId `bson:"favicon,omitempty" json:"favicon,omitempty"`

	// files
	Membership bson.ObjectId `bson:"membership,omitempty" json:"membership,omitempty"`

	PageSettings  map[string]*SitePageSettings  `bson:"page_settings" json:"-"`
	ThemeSettings map[string]*SiteThemeSettings `bson:"theme_settings" json:"-"`

	// build settings
	Theme        string `bson:"theme"         json:"theme"`
	Domain       string `bson:"domain"        json:"domain"`
	CustomDomain string `bson:"custom_domain" json:"customDomain"`
	CustomURL    string `bson:"custom_url"    json:"customUrl"`
	UglyURL      bool   `bson:"ugly_url"      json:"uglyUrl"`

	// theme settings
	NameInNavBar bool `bson:"name_in_navbar" json:"nameInNavBar"`
}

// SiteJSON represents the json version of a site
type SiteJSON struct {
	Site
	Links map[string]interface{} `json:"links"`

	// overrides the PageSettings and ThemeSettings fields to provide an array
	// of ids (as needed by Ember Data) instead of a hash of embedded documents
	PageSettings  []string `json:"pageSettings,omitempty"`
	ThemeSettings []string `json:"themeSettings,omitempty"`
}

// SitesList represents a list of sites
type SitesList []*Site

const (
	// PageKindContact represents the contact page
	PageKindContact = "contact"

	// PageKindActivities represents the activities page
	PageKindActivities = "activities"

	// PageKindMembers represents the members page
	PageKindMembers = "members"

	// PageKindPosts represents the contact page
	PageKindPosts = "posts"

	// PageKindEvents represents the contact page
	PageKindEvents = "events"
)

func init() {
	SitePagesSettingsKinds = map[string]bool{
		PageKindContact:    true,
		PageKindActivities: true,
		PageKindMembers:    true,
		PageKindPosts:      true,
		PageKindEvents:     true,
	}
}

//
// DBSession
//

// SitesCol returns the sites collection
func (session *DBSession) SitesCol() *mgo.Collection {
	return session.DB().C(sitesColName)
}

// EnsureSitesIndexes ensures indexes on sites collection
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

// FindSite finds a site by id
func (session *DBSession) FindSite(siteID string) *Site {
	var result Site

	if err := session.SitesCol().FindId(siteID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreateSite creates a new site in database
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

// RemoveImageReferencesFromSitePageSettings removes all references to given image from site page settings
func (session *DBSession) RemoveImageReferencesFromSitePageSettings(image *Image) error {
	// @todo
	return nil
}

//
// SiteThemeSettings
//

// MarshalJSON implements the json.Marshaler interface
func (settings *SiteThemeSettings) MarshalJSON() ([]byte, error) {
	customField, err := json.Marshal(settings.Custom)
	if err != nil {
		return []byte{}, err
	}

	settingsJSON := SiteThemeSettingsJSON{
		SiteThemeSettings: *settings,
	}

	customValue := string(customField)
	if customValue != "null" {
		settingsJSON.Custom = customValue
	}

	return json.Marshal(settingsJSON)
}

//
// Site
//

// MarshalJSON implements json.MarshalJSON
func (site *Site) MarshalJSON() ([]byte, error) {
	// inject 'links' needed by Ember Data
	// @todo Remove that ?
	links := map[string]interface{}{
		"posts":      fmt.Sprintf("/api/sites/%s/posts", site.ID),
		"events":     fmt.Sprintf("/api/sites/%s/events", site.ID),
		"pages":      fmt.Sprintf("/api/sites/%s/pages", site.ID),
		"activities": fmt.Sprintf("/api/sites/%s/activities", site.ID),
		"members":    fmt.Sprintf("/api/sites/%s/members", site.ID),
		"images":     fmt.Sprintf("/api/sites/%s/images", site.ID),
		"files":      fmt.Sprintf("/api/sites/%s/files", site.ID),
	}

	// convert hash of embedded docs into an array of doc ids, as needed by Ember Data
	pageSettingsIds := []string{}
	for _, settings := range site.PageSettings {
		pageSettingsIds = append(pageSettingsIds, settings.ID.Hex())
	}

	// convert hash of embedded docs into an array of doc ids, as needed by Ember Data
	themeSettingsIds := []string{}
	for _, settings := range site.ThemeSettings {
		themeSettingsIds = append(themeSettingsIds, settings.ID.Hex())
	}

	siteJSON := SiteJSON{
		Site:          *site,
		Links:         links,
		PageSettings:  pageSettingsIds,
		ThemeSettings: themeSettingsIds,
	}

	return json.Marshal(siteJSON)
}

// BaseUrl returns base URL for that site.
func (site *Site) BaseUrl() string {
	if site.CustomURL != "" {
		return site.CustomURL
	}

	if site.CustomDomain != "" {
		return core.BaseUrlForCustomDomain(site.CustomDomain)
	}

	if site.Domain != "" {
		return core.BaseUrlForDomain(site.ID, site.Domain)
	}

	return core.BaseUrl(site.ID)
}

// BuildDir returns build directory for that site.
func (site *Site) BuildDir() string {
	if site.CustomURL != "" {
		u, err := url.Parse(site.CustomURL)
		if err == nil {
			if helpers.HasOnePrefix(u.Host, []string{"127.0.0.1", "localhost"}) || u.Host == "" {
				// eg: 3ailes
				return site.ID
			}

			// eg: 3ailes.org
			return u.Host
		}
	}

	if site.CustomDomain != "" {
		// eg: 3ailes.org
		return site.CustomDomain
	}

	if site.Domain != "" {
		// eg: 3ailes.asso.ninja
		return site.ID + "." + site.Domain
	}

	// eg: 3ailes
	return site.ID
}

// TZLocation returns timezone Location
func (site *Site) TZLocation() *time.Location {
	tz := site.TZ
	if tz == "" {
		tz = core.DefaultTZ
	}

	result, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}

	return result
}

//
// Site posts
//

func (site *Site) postsBaseQuery(onlyPub bool) *mgo.Query {
	selector := bson.M{"site_id": site.ID}

	if onlyPub {
		selector["published"] = true
	}

	return site.dbSession.PostsCol().Find(selector)
}

// PostsNb returns the total number of posts
func (site *Site) PostsNb() int {
	result, err := site.postsBaseQuery(false).Count()
	if err != nil {
		panic(err)
	}

	return result
}

// FindPosts fetches posts belonging to site
func (site *Site) FindPosts(skip int, limit int, onlyPub bool) *PostsList {
	result := PostsList{}

	query := site.postsBaseQuery(onlyPub).Sort("published", "-published_at", "-updated_at")

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

// FindAllPosts fetches all posts belonging to site
func (site *Site) FindAllPosts() *PostsList {
	return site.FindPosts(0, 0, false)
}

// FindPublishedPosts fetches all published posts belonging to site
func (site *Site) FindPublishedPosts() *PostsList {
	return site.FindPosts(0, 0, true)
}

//
// Site events
//

func (site *Site) eventsBaseQuery() *mgo.Query {
	return site.dbSession.EventsCol().Find(bson.M{"site_id": site.ID})
}

// EventsNb returns the total number of events
func (site *Site) EventsNb() int {
	result, err := site.eventsBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// FindEvents fetches events belonging to site
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

// FindAllEvents fetches all events belonging to site
func (site *Site) FindAllEvents() *EventsList {
	return site.FindEvents(0, 0)
}

//
// Site pages
//

func (site *Site) pagesBaseQuery() *mgo.Query {
	return site.dbSession.PagesCol().Find(bson.M{"site_id": site.ID})
}

// PagesNb returns the total number of pages
func (site *Site) PagesNb() int {
	result, err := site.pagesBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// FindPages fetches pages belonging to site
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

// FindAllPages fetches all pages belonging to site
func (site *Site) FindAllPages() *PagesList {
	return site.FindPages(0, 0)
}

//
// Site activities
//

func (site *Site) activitiesBaseQuery() *mgo.Query {
	return site.dbSession.ActivitiesCol().Find(bson.M{"site_id": site.ID})
}

// ActivitiesNb returns the total number of activities
func (site *Site) ActivitiesNb() int {
	result, err := site.activitiesBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// FindActivities fetches activities belonging to site
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

// FindAllActivities fetches all activities belonging to site
func (site *Site) FindAllActivities() *ActivitiesList {
	return site.FindActivities(0, 0)
}

//
// Site members
//

func (site *Site) membersBaseQuery() *mgo.Query {
	return site.dbSession.MembersCol().Find(bson.M{"site_id": site.ID})
}

// MembersNb returns the total number of members
func (site *Site) MembersNb() int {
	result, err := site.membersBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// FindMembers fetches members belonging to site
func (site *Site) FindMembers(skip int, limit int) *MembersList {
	result := MembersList{}

	query := site.membersBaseQuery().Sort("order", "created_at")

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

// FindAllMembers fetches all members belonging to site
func (site *Site) FindAllMembers() *MembersList {
	return site.FindMembers(0, 0)
}

// UpdateMemberOrder updates a member order in database
func (site *Site) UpdateMemberOrder(id bson.ObjectId, order int) {
	// specify site_id in selector to prevent unprivileged users access
	selector := bson.M{"site_id": site.ID, "_id": id}
	modifier := bson.M{"$set": bson.M{"order": order}}

	site.dbSession.MembersCol().Update(selector, modifier)
}

//
// Site images
//

func (site *Site) imagesBaseQuery() *mgo.Query {
	return site.dbSession.ImagesCol().Find(bson.M{"site_id": site.ID})
}

// ImagesNb returns the total number of images
func (site *Site) ImagesNb() int {
	result, err := site.imagesBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// FindImages fetches images belonging to site
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

// FindAllImages fetches all images belonging to site
func (site *Site) FindAllImages() *ImagesList {
	return site.FindImages(0, 0)
}

//
// Site files
//

func (site *Site) filesBaseQuery() *mgo.Query {
	return site.dbSession.FilesCol().Find(bson.M{"site_id": site.ID})
}

// FilesNb returns the total number of files
func (site *Site) FilesNb() int {
	result, err := site.filesBaseQuery().Count()
	if err != nil {
		panic(err)
	}

	return result
}

// FindFiles fetches files belonging to site
func (site *Site) FindFiles(skip int, limit int) *FilesList {
	result := FilesList{}

	query := site.filesBaseQuery().Sort("-created_at")

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

// FindAllFiles fetches all files belonging to site
func (site *Site) FindAllFiles() *FilesList {
	return site.FindFiles(0, 0)
}

//
// Site fields
//

// FindLogo fetches logo from database
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

// FindCover fetches cover from database
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

// FindFavicon fetches favicon from database
func (site *Site) FindFavicon() *Image {
	if site.Favicon != "" {
		var result Image

		if err := site.dbSession.ImagesCol().FindId(site.Favicon).One(&result); err != nil {
			return nil
		}

		result.dbSession = site.dbSession

		return &result
	}

	return nil
}

// FindMembership fetches membership from database
func (site *Site) FindMembership() *File {
	if site.Membership != "" {
		var result File

		if err := site.dbSession.FilesCol().FindId(site.Membership).One(&result); err != nil {
			return nil
		}

		result.dbSession = site.dbSession

		return &result
	}

	return nil
}

// FindPageSettingsCover fetches page settings cover from database
func (site *Site) FindPageSettingsCover(settingKind string) *Image {
	pageSettings := site.PageSettings[settingKind]
	if (pageSettings != nil) && (pageSettings.Cover != "") {
		var result Image

		if err := site.dbSession.ImagesCol().FindId(pageSettings.Cover).One(&result); err != nil {
			return nil
		}

		result.dbSession = site.dbSession

		return &result
	}

	return nil
}

// RemoveImageReferences removes all references to given image from database
func (site *Site) RemoveImageReferences(image *Image) error {
	// remove image reference from site settings
	fieldsToDelete := []string{}

	if site.Logo == image.ID {
		fieldsToDelete = append(fieldsToDelete, "logo")
	}

	if site.Cover == image.ID {
		fieldsToDelete = append(fieldsToDelete, "cover")
	}

	if site.Favicon == image.ID {
		fieldsToDelete = append(fieldsToDelete, "favicon")
	}

	if len(fieldsToDelete) > 0 {
		site.DeleteFields(fieldsToDelete)
	}

	// remove image references from page settings
	if err := site.dbSession.RemoveImageReferencesFromSitePageSettings(image); err != nil {
		return err
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

// RemoveFileReferences removes all references to given file from database
func (site *Site) RemoveFileReferences(file *File) error {
	// remove image reference from site settings
	fieldsToDelete := []string{}

	if site.Membership == file.ID {
		fieldsToDelete = append(fieldsToDelete, "membership")
	}

	if len(fieldsToDelete) > 0 {
		site.DeleteFields(fieldsToDelete)
	}

	return nil
}

// Update updates site in database
func (site *Site) Update(newSite *Site) (bool, error) {
	var set, unset, modifier bson.D

	if site.Lang != newSite.Lang {
		site.Lang = newSite.Lang

		if site.Lang == "" {
			unset = append(unset, bson.DocElem{"lang", 1})
		} else {
			set = append(set, bson.DocElem{"lang", site.Lang})
		}
	}

	if site.TZ != newSite.TZ {
		site.TZ = newSite.TZ

		if site.TZ == "" {
			unset = append(unset, bson.DocElem{"tz", 1})
		} else {
			set = append(set, bson.DocElem{"tz", site.TZ})
		}
	}

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

	if site.GoogleAnalytics != newSite.GoogleAnalytics {
		site.GoogleAnalytics = newSite.GoogleAnalytics

		if site.GoogleAnalytics == "" {
			unset = append(unset, bson.DocElem{"google_analytics", 1})
		} else {
			set = append(set, bson.DocElem{"google_analytics", site.GoogleAnalytics})
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

	if site.Favicon != newSite.Favicon {
		site.Favicon = newSite.Favicon

		if site.Favicon == "" {
			unset = append(unset, bson.DocElem{"favicon", 1})
		} else {
			set = append(set, bson.DocElem{"favicon", site.Favicon})
		}
	}

	if site.Membership != newSite.Membership {
		site.Membership = newSite.Membership

		if site.Membership == "" {
			unset = append(unset, bson.DocElem{"membership", 1})
		} else {
			set = append(set, bson.DocElem{"membership", site.Membership})
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

	if site.Domain != newSite.Domain {
		site.Domain = newSite.Domain

		if site.Domain == "" {
			unset = append(unset, bson.DocElem{"domain", 1})
		} else {
			set = append(set, bson.DocElem{"domain", site.Domain})
		}
	}

	if site.CustomDomain != newSite.CustomDomain {
		site.CustomDomain = newSite.CustomDomain

		if site.CustomDomain == "" {
			unset = append(unset, bson.DocElem{"custom_domain", 1})
		} else {
			set = append(set, bson.DocElem{"custom_domain", site.CustomDomain})
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

	if site.NameInNavBar != newSite.NameInNavBar {
		site.NameInNavBar = newSite.NameInNavBar

		if site.NameInNavBar == false {
			unset = append(unset, bson.DocElem{"name_in_navbar", 1})
		} else {
			set = append(set, bson.DocElem{"name_in_navbar", site.NameInNavBar})
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
		return true, site.dbSession.SitesCol().UpdateId(site.ID, modifier)
	}

	return false, nil
}

// SetValues sets given site fields in database
func (site *Site) SetValues(values bson.M) error {
	// @todo Set UpdatedAt field
	return site.dbSession.SitesCol().UpdateId(site.ID, bson.D{{"$set", values}})
}

// SetChangedAt sets the ChangedAt value
func (site *Site) SetChangedAt(value time.Time) error {
	if err := site.SetValues(bson.M{"changed_at": value}); err != nil {
		return err
	}

	site.ChangedAt = value
	return nil
}

// SetBuiltAt sets the BuiltAt value
func (site *Site) SetBuiltAt(value time.Time) error {
	if err := site.SetValues(bson.M{"built_at": value}); err != nil {
		return err
	}

	site.BuiltAt = value
	return nil
}

// SetMembership sets the Membership value
func (site *Site) SetMembership(value bson.ObjectId) error {
	if err := site.SetValues(bson.M{"membership": value}); err != nil {
		return err
	}

	site.Membership = value
	return nil
}

// DeleteFields delete given site fields from database
func (site *Site) DeleteFields(fields []string) error {
	var unset bson.D

	for _, field := range fields {
		unset = append(unset, bson.DocElem{field, 1})
	}

	return site.dbSession.SitesCol().UpdateId(site.ID, bson.M{"$unset": unset})
}

// Delete site from database
func (site *Site) Delete() error {
	// delete site
	if err := site.dbSession.SitesCol().RemoveId(site.ID); err != nil {
		return err
	}

	// delete site content
	// @todo Catch and report errors
	site.dbSession.ActivitiesCol().RemoveAll(bson.M{"site_id": site.ID})
	site.dbSession.EventsCol().RemoveAll(bson.M{"site_id": site.ID})
	site.dbSession.ImagesCol().RemoveAll(bson.M{"site_id": site.ID})
	site.dbSession.MembersCol().RemoveAll(bson.M{"site_id": site.ID})
	site.dbSession.PagesCol().RemoveAll(bson.M{"site_id": site.ID})
	site.dbSession.PostsCol().RemoveAll(bson.M{"site_id": site.ID})

	// delete site images
	// @todo Catch and report error
	site.deleteImagesFiles()

	return nil
}

// Delete all images files
func (site *Site) deleteImagesFiles() error {
	dirPath := core.UploadSiteDir(site.ID)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		if errRem := os.RemoveAll(dirPath); errRem != nil {
			return errRem
		}
	}

	return nil
}

// SetPageSettings inserts (or updates) page settings to database
// Side effect: 'Id' field is set on record if not already present, and string fields are trimed
func (site *Site) SetPageSettings(settings *SitePageSettings) error {
	if !SitePagesSettingsKinds[settings.Kind] {
		return errors.New("Unsupported page settings kind: " + settings.Kind)
	}

	if settings.ID == "" {
		settings.ID = bson.NewObjectId()
	}

	settings.Title = strings.TrimSpace(settings.Title)
	settings.Tagline = strings.TrimSpace(settings.Tagline)

	return site.dbSession.SitesCol().UpdateId(site.ID, bson.M{"$set": bson.D{bson.DocElem{fmt.Sprintf("page_settings.%s", settings.Kind), settings}}})
}

// SetThemeSettings inserts (or updates) theme settings to database
// Side effect: 'Id' field is set on record if not already present, and string fields are trimed
func (site *Site) SetThemeSettings(settings *SiteThemeSettings) error {
	if !themes.Exist(settings.Theme) {
		return errors.New("Theme does not exist: " + settings.Theme)
	}

	if settings.ID == "" {
		settings.ID = bson.NewObjectId()
	}

	return site.dbSession.SitesCol().UpdateId(site.ID, bson.M{"$set": bson.D{bson.DocElem{fmt.Sprintf("theme_settings.%s", settings.Theme), settings}}})
}
