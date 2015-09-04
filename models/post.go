package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	postsColName = "posts"
)

// Post represents a post
type Post struct {
	dbSession *DBSession `bson:"-"`

	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteID    string        `bson:"site_id"       json:"site"`

	Published   bool          `bson:"published"       json:"published"`
	PublishedAt time.Time     `bson:"published_at"    json:"publishedAt,omitempty"`
	Title       string        `bson:"title"           json:"title"`
	Body        string        `bson:"body"            json:"body"`
	Format      string        `bson:"format"          json:"format"`
	Cover       bson.ObjectId `bson:"cover,omitempty" json:"cover,omitempty"`
}

// PostsList represents a list of posts
type PostsList []*Post

//
// DBSession
//

// PostsCol returns posts collection
func (session *DBSession) PostsCol() *mgo.Collection {
	return session.DB().C(postsColName)
}

// EnsurePostsIndexes ensures indexes on posts collection
func (session *DBSession) EnsurePostsIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id", "published"},
		Background: true,
	}

	err := session.PostsCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// FindPost finds a post by id
func (session *DBSession) FindPost(postID bson.ObjectId) *Post {
	var result Post

	if err := session.PostsCol().FindId(postID).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// CreatePost creates a new post in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on post record
func (session *DBSession) CreatePost(post *Post) error {
	post.ID = bson.NewObjectId()

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	if err := session.PostsCol().Insert(post); err != nil {
		return err
	}

	post.dbSession = session

	return nil
}

// RemoveImageReferencesFromPosts removes all references to given image from all posts
func (session *DBSession) RemoveImageReferencesFromPosts(image *Image) error {
	// @todo
	return nil
}

//
// Post
//

// FindSite fetches site that post belongs to
func (post *Post) FindSite() *Site {
	return post.dbSession.FindSite(post.SiteID)
}

// FindCover fetches cover from database
func (post *Post) FindCover() *Image {
	if post.Cover != "" {
		var result Image

		if err := post.dbSession.ImagesCol().FindId(post.Cover).One(&result); err != nil {
			return nil
		}

		result.dbSession = post.dbSession

		return &result
	}

	return nil
}

// Delete deletes post from database
func (post *Post) Delete() error {
	// delete from database
	if err := post.dbSession.PostsCol().RemoveId(post.ID); err != nil {
		return err
	}

	return nil
}

// Update updates post in database
func (post *Post) Update(newPost *Post) (bool, error) {
	var set, unset, modifier bson.D

	// Published
	if post.Published != newPost.Published {
		post.Published = newPost.Published

		set = append(set, bson.DocElem{"published", post.Published})

		// PublishedAt
		if post.Published {
			post.PublishedAt = time.Now()

			set = append(set, bson.DocElem{"published_at", post.PublishedAt})
		}
	} else if post.Published && (post.PublishedAt != newPost.PublishedAt) {
		// PublishedAt
		post.PublishedAt = newPost.PublishedAt

		set = append(set, bson.DocElem{"published_at", post.PublishedAt})
	}

	// Title
	if post.Title != newPost.Title {
		post.Title = newPost.Title

		if post.Title == "" {
			unset = append(unset, bson.DocElem{"title", 1})
		} else {
			set = append(set, bson.DocElem{"title", post.Title})
		}
	}

	// Body
	if post.Body != newPost.Body {
		post.Body = newPost.Body

		if post.Body == "" {
			unset = append(unset, bson.DocElem{"body", 1})
		} else {
			set = append(set, bson.DocElem{"body", post.Body})
		}
	}

	// Format
	newFormat := newPost.Format
	if newFormat == "" {
		newFormat = DefaultFormat
	}

	if post.Format != newFormat {
		post.Format = newFormat

		set = append(set, bson.DocElem{"format", post.Format})
	}

	// Cover
	if post.Cover != newPost.Cover {
		post.Cover = newPost.Cover

		if post.Cover == "" {
			unset = append(unset, bson.DocElem{"cover", 1})
		} else {
			set = append(set, bson.DocElem{"cover", post.Cover})
		}
	}

	if len(unset) > 0 {
		modifier = append(modifier, bson.DocElem{"$unset", unset})
	}

	if len(set) > 0 {
		modifier = append(modifier, bson.DocElem{"$set", set})
	}

	if len(modifier) > 0 {
		post.UpdatedAt = time.Now()
		set = append(set, bson.DocElem{"updated_at", post.UpdatedAt})

		return true, post.dbSession.PostsCol().UpdateId(post.ID, modifier)
	}

	return false, nil
}
