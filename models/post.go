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
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
	SiteId    bson.ObjectId `bson:"site_id"       json:"site"`

	PublishedAt time.Time `bson:"published_at" json:"publishedAt"`
	Title       string    `bson:"title"        json:"title"`
	Body        string    `bson:"body"         json:"body"`
	// @todo Photo
}

type PostsList []Post

// Posts collection
func PostsCol() *mgo.Collection {
	return DB().C(POSTS_COL_NAME)
}
