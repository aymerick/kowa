package builder

import (
	"fmt"
	"html/template"
	"path"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/aymerick/kowa/models"
)

// Builder for posts pages
type PostsBuilder struct {
	*NodeBuilderBase

	posts []*PostNodeContentPair
}

// Post content for template
type PostContent struct {
	Node *Node

	Date  string
	Cover string
	Title string
	Body  template.HTML
	Url   string
}

// Post with associated Node Content
type PostNodeContentPair struct {
	post        *models.Post
	nodeContent *PostContent
}

// Post list content for template
type PostListContent struct {
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

func NewPostNodeContentPair(post *models.Post, nodeContent *PostContent) *PostNodeContentPair {
	return &PostNodeContentPair{
		post:        post,
		nodeContent: nodeContent,
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

// Build post page
func (builder *PostsBuilder) loadPost(post *models.Post) {
	node := builder.newNode()
	node.fillUrl(path.Join("posts", post.Slug())) // @todo i18n

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

	builder.posts = append(builder.posts, NewPostNodeContentPair(post, postContent))
}

// Instanciate a new post content
func (builder *PostsBuilder) NewPostContent(post *models.Post, node *Node) *PostContent {
	result := &PostContent{
		Node:  node,
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

		node.Content = &PostListContent{
			Node:  node,
			Posts: computesPostContents(builder.posts),
		}

		builder.addNode(node)
	}
}

func computesPostContents(posts []*PostNodeContentPair) []*PostContent {
	postContents := []*PostContent{}

	for _, postNodeContent := range posts {
		postContents = append(postContents, postNodeContent.nodeContent)
	}

	return postContents
}
