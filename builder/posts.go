package builder

import (
	"html/template"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Post content for template
type PostContent struct {
	Date    time.Time     // CreatedAt
	Image   string        // Cover
	Title   string        // Title
	Content template.HTML // Body
	Url     string        // Absolute URL
}

// Builder for posts pages
type PostsBuilder struct {
	*NodeBuilder
}

// Instanciate a new builder
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

// Build posts list pages
func (builder *PostsBuilder) BuildPostsLists() {
	// @todo !!!
}

// Build posts single pages
func (builder *PostsBuilder) BuildPosts() {
	for _, post := range *builder.Site().Model.FindAllPosts() {
		builder.BuildPost(post)
	}
}

// Build post single page
func (builder *PostsBuilder) BuildPost(post *models.Post) {
	node := builder.NewNode()

	node.basePath = utils.Urlify(post.Title)

	node.Title = post.Title
	node.Meta = &NodeMeta{
		Description: "@todo",
	}

	node.Content = builder.NewPostContent(post, node)

	builder.AddNode(node)
}

/// Instanciate a new post content
func (builder *PostsBuilder) NewPostContent(post *models.Post, node *Node) *PostContent {
	result := &PostContent{
		Date:  post.CreatedAt,
		Title: post.Title,
		Url:   node.Url(),
	}

	cover := post.FindCover()
	if cover != nil {
		result.Image = cover.MediumURL()
	}

	html := blackfriday.MarkdownCommon([]byte(post.Body))
	result.Content = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	return result
}
