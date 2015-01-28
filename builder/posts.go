package builder

import (
	"html/template"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Builder for posts pages
type PostsBuilder struct {
	*NodeBuilderBase
}

// Post content for template
type PostContent struct {
	Date  time.Time     // CreatedAt
	Cover string        // Cover
	Title string        // Title
	Body  template.HTML // Body
	Url   string        // Absolute URL
}

func init() {
	RegisterNodeBuilder(KIND_POSTS, NewPostsBuilder)
}

// Instanciate a new builder
func NewPostsBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &PostsBuilder{
		&NodeBuilderBase{
			nodeKind:    KIND_POST,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *PostsBuilder) Load() {
	builder.loadPosts()
	builder.loadPostsLists()
}

// Build all posts
func (builder *PostsBuilder) loadPosts() {
	for _, post := range *builder.SiteBuilder().site.FindAllPosts() {
		builder.loadPost(post)
	}
}

// Build post page
func (builder *PostsBuilder) loadPost(post *models.Post) {
	node := builder.newNode()
	node.fillUrl(post.Slug())

	node.Title = post.Title
	node.Meta = &NodeMeta{
		Description: "@todo",
	}

	node.Content = builder.NewPostContent(post, node)

	builder.addNode(node)
}

/// Instanciate a new post content
func (builder *PostsBuilder) NewPostContent(post *models.Post, node *Node) *PostContent {
	result := &PostContent{
		Date:  post.CreatedAt,
		Title: post.Title,
		Url:   node.Url,
	}

	cover := post.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover, models.MEDIUM_KIND)
	}

	html := blackfriday.MarkdownCommon([]byte(post.Body))
	result.Body = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	return result
}

// Build posts list pages
func (builder *PostsBuilder) loadPostsLists() {
	node := builder.newNodeForKind(KIND_POSTS)
	node.fillUrl(KIND_POSTS)

	node.Title = "Posts"
	node.Meta = &NodeMeta{Description: "Posts test node"}
	node.Content = "Soon"
	node.InNavBar = true

	builder.addNode(node)
}
