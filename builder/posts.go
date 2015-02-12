package builder

import (
	"fmt"
	"html/template"
	"path"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
)

// Post nodes builder
type PostsBuilder struct {
	*NodeBuilderBase

	posts []*PostContent
}

// Post node content
type PostContent struct {
	Node  *Node
	Model *models.Post

	Date  string
	Cover string
	Title string
	Body  template.HTML
	Url   string
}

// Posts node content
type PostsContent struct {
	Node *Node

	Posts []*PostContent
	// PrevPage string
	// NextPage string
}

func init() {
	RegisterNodeBuilder(KIND_POSTS, NewPostsBuilder)
}

// Instanciate a new builder
func NewPostsBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &PostsBuilder{
		NodeBuilderBase: &NodeBuilderBase{
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
	for _, post := range *builder.site().FindAllPosts() {
		builder.loadPost(post)
	}
}

// Computes slug
func postSlug(post *models.Post) string {
	// @todo Should use PublishedAt
	year, month, day := post.CreatedAt.Date()

	title := post.Title
	if len(title) > MAX_SLUG {
		title = title[:MAX_SLUG]
	}

	return fmt.Sprintf("%d/%02d/%02d/%s", year, month, day, utils.Urlify(title))
}

// Build post page
func (builder *PostsBuilder) loadPost(post *models.Post) {
	node := builder.newNode()
	node.fillUrl(path.Join("posts", postSlug(post))) // @todo i18n

	title := "Posts" // @todo i18n
	tagline := ""    // @todo fill

	node.Title = title
	node.Tagline = tagline
	node.Meta = &NodeMeta{
		Title:       fmt.Sprintf("%s - %s", post.Title, builder.site().Name),
		Description: tagline,
	}

	postContent := builder.NewPostContent(post, node)
	node.Content = postContent

	builder.addNode(node)

	builder.posts = append(builder.posts, postContent)
}

// Instanciate a new post content
func (builder *PostsBuilder) NewPostContent(post *models.Post, node *Node) *PostContent {
	result := &PostContent{
		Node:  node,
		Model: post,

		Date:  post.CreatedAt.Format("02/01/06"),
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
	if len(builder.posts) > 0 {
		// @todo pagination
		node := builder.newNodeForKind(KIND_POSTS)
		node.fillUrl(KIND_POSTS)

		title := "Posts" // @todo i18n
		tagline := ""    // @todo fill

		node.Title = title
		node.Tagline = tagline
		node.Meta = &NodeMeta{Description: tagline}
		node.InNavBar = true
		node.NavBarOrder = 5

		node.Content = &PostsContent{
			Node:  node,
			Posts: builder.posts,
		}

		builder.addNode(node)
	}
}
