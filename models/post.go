package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	POSTS_COL_NAME = "posts"
)

type Post struct {
	dbSession *DBSession `bson:"-" json:"-"`

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	PublishedAt *time.Time `bson:"published_at" json:"publishedAt,omitempty"`
	Title       string     `bson:"title"        json:"title"`
	Body        string     `bson:"body"         json:"body"`
	// @todo Photo/Cover
	// @todo Format (markdown|html)
}

type PostsList []*Post

//
// DBSession
//

// Posts collection
func (session *DBSession) PostsCol() *mgo.Collection {
	return session.DB().C(POSTS_COL_NAME)
}

// Ensure indexes on Posts collection
func (session *DBSession) EnsurePostsIndexes() {
	index := mgo.Index{
		Key:        []string{"site_id"},
		Background: true,
	}

	err := session.PostsCol().EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// Find post by id
func (session *DBSession) FindPost(postId bson.ObjectId) *Post {
	var result Post

	if err := session.PostsCol().FindId(postId).One(&result); err != nil {
		return nil
	}

	result.dbSession = session

	return &result
}

// Persists a new post in database
// Side effect: 'Id', 'CreatedAt' and 'UpdatedAt' fields are set on post record
func (session *DBSession) CreatePost(post *Post) error {
	post.Id = bson.NewObjectId()

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	if err := session.PostsCol().Insert(post); err != nil {
		return err
	}

	post.dbSession = session

	return nil
}

//
// Post
//

// Fetch from database: site that post belongs to
func (post *Post) FindSite() *Site {
	return post.dbSession.FindSite(post.SiteId)
}

// Delete post from database
func (post *Post) Delete() error {
	var err error

	// delete from database
	if err = post.dbSession.PostsCol().RemoveId(post.Id); err != nil {
		return err
	}

	return nil
}

// Update post in database
func (post *Post) Update(newPost *Post) error {
	var set, unset, modifier bson.D

	if post.Title != newPost.Title {
		post.Title = newPost.Title

		if post.Title == "" {
			unset = append(unset, bson.DocElem{"title", 1})
		} else {
			set = append(set, bson.DocElem{"title", post.Title})
		}
	}

	if post.Body != newPost.Body {
		post.Body = newPost.Body

		if post.Body == "" {
			unset = append(unset, bson.DocElem{"body", 1})
		} else {
			set = append(set, bson.DocElem{"body", post.Body})
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

		return post.dbSession.PostsCol().UpdateId(post.Id, modifier)
	} else {
		return nil
	}
}
