package builder

import (
	"fmt"
	"path"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/raymond"
)

// PostsBuilder builds posts
type PostsBuilder struct {
	*NodeBuilderBase

	posts []*PostContent
}

// PostContent represents a post node content
type PostContent struct {
	Model *models.Post

	Date  string
	Cover *ImageVars
	Title string
	Body  raymond.SafeString
	Url   string
}

// PostsContent represents the posts page content
type PostsContent struct {
	Posts []*PostContent
	// PrevPage string
	// NextPage string
}

func init() {
	RegisterNodeBuilder(kindPosts, NewPostsBuilder)
}

// NewPostsBuilder instanciate a new NodeBuilder
func NewPostsBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &PostsBuilder{
		NodeBuilderBase: &NodeBuilderBase{
			nodeKind:    kindPost,
			siteBuilder: siteBuilder,
		},
	}
}

// Load is part of NodeBuilder interface
func (builder *PostsBuilder) Load() {
	builder.loadPosts()
	builder.loadPostsLists()
}

// Build published posts
func (builder *PostsBuilder) loadPosts() {
	for _, post := range *builder.site().FindPublishedPosts() {
		builder.loadPost(post)
	}
}

// Computes slug
func postSlug(post *models.Post) string {
	year, month, day := post.PublishedAt.Date()

	title := post.Title
	if len(title) > maxSlug {
		title = title[:maxSlug]
	}

	return fmt.Sprintf("%d/%02d/%02d/%s", year, month, day, title)
}

// Build post page
func (builder *PostsBuilder) loadPost(post *models.Post) {
	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PAGE_KIND_POSTS)
	if disabled {
		return
	}

	T := i18n.MustTfunc(builder.siteLang())
	slug := T("posts")

	if title == "" {
		title = slug
	}

	// build node
	node := builder.newNode()
	node.fillURL(path.Join(slug, postSlug(post)))

	node.Title = title
	node.Tagline = tagline

	node.Meta = &NodeMeta{
		Title:       fmt.Sprintf("%s - %s", post.Title, builder.site().Name),
		Description: tagline,
		Type:        "article",
	}

	postContent := builder.NewPostContent(post, node)
	node.Content = postContent

	if postContent.Cover != nil {
		node.Cover = postContent.Cover
	} else {
		node.Cover = cover
	}

	builder.addNode(node)

	builder.posts = append(builder.posts, postContent)
}

// NewPostContent instanciate a new PostContent
func (builder *PostsBuilder) NewPostContent(post *models.Post, node *Node) *PostContent {
	T := i18n.MustTfunc(builder.siteLang())

	result := &PostContent{
		Model: post,

		Title: post.Title,
		Url:   node.Url,
	}

	year, _, day := post.PublishedAt.Date()
	result.Date = T("post_format_date", map[string]interface{}{
		"Year":  year,
		"Month": T("month_" + post.PublishedAt.Format("January")),
		"Day":   day,
	})

	cover := post.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover)
	}

	result.Body = generateHTML(post.Format, post.Body)

	return result
}

// Build posts list pages
// @todo pagination
func (builder *PostsBuilder) loadPostsLists() {
	if len(builder.posts) == 0 {
		return
	}

	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PAGE_KIND_POSTS)
	if disabled {
		return
	}

	T := i18n.MustTfunc(builder.siteLang())
	slug := T("posts")

	if title == "" {
		title = slug
	}

	// build node
	node := builder.newNodeForKind(kindPosts)
	node.fillURL(slug)

	node.Title = title
	node.Tagline = tagline
	node.Cover = cover

	node.Meta = &NodeMeta{Description: tagline}
	node.InNavBar = true
	node.NavBarOrder = 5

	node.Content = &PostsContent{
		Posts: builder.posts,
	}

	builder.addNode(node)
}
