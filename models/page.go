package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	pagesColName = "pages"
)

// Page represents a page
type Page struct {
	dbSession *DBSession `bson:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	Title   string        `bson:"title"           json:"title"`
	Tagline string        `bson:"tagline"         json:"tagline"`
	Body    string        `bson:"body"            json:"body"`
	Format  string        `bson:"format"          json:"format"`
	Cover   bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`

	InNavBar bool `bson:"in_nav_bar" json:"inNavBar"`
}

// PagesList represents a list of pages
type PagesList []*Page

//
// DBSession
//

// PagesCol returns pages collection
func (session *DBSession) PagesCol() *mgo.Collection {
	return session.DB().C(pagesColName)
}

// EnsurePagesIndexes ensure indexes on pages collection
func (session *DBSession) EnsurePagesIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.PagesCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// FindPage finds a page by id
func (session *DBSession) FindPage(pageID bson.ObjectId) *Page {
	var result Page

	if err := session.PagesCol().FindId(pageID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreatePage creates a new page in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on page record
func (session *DBSession) CreatePage(page *Page) error {
	page.Id = bson.NewObjectId()

	now := time.Now()
	page.CreatedAt = now
	page.UpdatedAt = now

	if err := session.PagesCol().Insert(page); err != nil {
		return err
	}

	page.dbSession = session

	return nil
}

// RemoveImageReferencesFromPages removes all references to given image from all pages
func (session *DBSession) RemoveImageReferencesFromPages(image *Image) error {
	// @todo
	return nil
}

//
// Page
//

// FindSite fetches site that page belongs to
func (page *Page) FindSite() *Site {
	return page.dbSession.FindSite(page.SiteId)
}

// FindCover fetches cover from database
func (page *Page) FindCover() *Image {
	if page.Cover != "" {
		var result Image

		if err := page.dbSession.ImagesCol().FindId(page.Cover).One(&result); err != nil {
			return nil
		}

		result.dbSession = page.dbSession

		return &result
	}

	return nil
}

// Delete page from database
func (page *Page) Delete() error {
	var err error

	// delete from database
	if err = page.dbSession.PagesCol().RemoveId(page.Id); err != nil {
		return err
	}

	return nil
}

// Update page in database
func (page *Page) Update(newPage *Page) (bool, error) {
	var set, unset, modifier bson.D

	// Title
	if page.Title != newPage.Title {
		page.Title = newPage.Title

		if page.Title == "" {
			unset = append(unset, bson.DocElem{"title", 1})
		} else {
			set = append(set, bson.DocElem{"title", page.Title})
		}
	}

	// Tagline
	if page.Tagline != newPage.Tagline {
		page.Tagline = newPage.Tagline

		if page.Tagline == "" {
			unset = append(unset, bson.DocElem{"tagline", 1})
		} else {
			set = append(set, bson.DocElem{"tagline", page.Tagline})
		}
	}

	// Body
	if page.Body != newPage.Body {
		page.Body = newPage.Body

		if page.Body == "" {
			unset = append(unset, bson.DocElem{"body", 1})
		} else {
			set = append(set, bson.DocElem{"body", page.Body})
		}
	}

	// Format
	newFormat := newPage.Format
	if newFormat == "" {
		newFormat = DefaultFormat
	}

	if page.Format != newFormat {
		page.Format = newFormat

		set = append(set, bson.DocElem{"format", page.Format})
	}

	// Cover
	if page.Cover != newPage.Cover {
		page.Cover = newPage.Cover

		if page.Cover == "" {
			unset = append(unset, bson.DocElem{"cover", 1})
		} else {
			set = append(set, bson.DocElem{"cover", page.Cover})
		}
	}

	// InNavBar
	if page.InNavBar != newPage.InNavBar {
		page.InNavBar = newPage.InNavBar

		if page.InNavBar == false {
			unset = append(unset, bson.DocElem{"in_nav_bar", 1})
		} else {
			set = append(set, bson.DocElem{"in_nav_bar", page.InNavBar})
		}
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		page.UpdatedAt = time.Now()
		set = append(set, bson.DocElem{"updated_at", page.UpdatedAt})

		return true, page.dbSession.PagesCol().UpdateId(page.Id, modifier)
	}

	return false, nil
}
