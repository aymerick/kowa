package builder

import (
	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
)

// Builder for posts pages
type PostsBuilder struct {
	*NodeBuilder
}

func NewPostsBuilder(site *Site) *PostsBuilder {
	return &PostsBuilder{
		&NodeBuilder{
			NodeKind: KIND_POST,
			site:     site,
		},
	}
}

// NodeBuilderInterface
func (builder *PostsBuilder) Load() {
	builder.BuildPostsLists()
	builder.BuildPosts()
}

func (builder *PostsBuilder) BuildPostsLists() {
	// @todo !!!
}

func (builder *PostsBuilder) BuildPosts() {
	for _, post := range *builder.Site().Model.FindAllPosts() {
		builder.BuildPost(post)
	}
}

func (builder *PostsBuilder) BuildPost(post *models.Post) {
	node := builder.NewNode()

	node.basePath = utils.Urlify(post.Title)

	node.Title = post.Title

	node.Meta = &NodeMeta{
		Description: "@todo",
	}

	node.Content = post

	builder.AddNode(node)
}
