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

	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    string        `bson:"site_id"       json:"site"`

	PublishedAt time.Time `bson:"published_at" json:"publishedAt"`
	Title       string    `bson:"title"        json:"title"`
	Body        string    `bson:"body"         json:"body"`
	// @todo Photo
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

//
// Post
//

// Fetch from database: site that post belongs to
func (post *Post) FindSite() *Site {
	return post.dbSession.FindSite(post.SiteId)
}
