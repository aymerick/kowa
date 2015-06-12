package builder

import (
	"fmt"
	"path"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/raymond"
)

// Post nodes builder
type PostsBuilder struct {
	*NodeBuilderBase

	posts []*PostContent
}

// Post node content
type PostContent struct {
	Model *models.Post

	Date  string
	Cover *ImageVars
	Title string
	Body  raymond.SafeString
	Url   string
}

// Posts node content
type PostsContent struct {
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
	node.fillUrl(path.Join(slug, postSlug(post)))

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

// Instanciate a new post content
func (builder *PostsBuilder) NewPostContent(post *models.Post, node *Node) *PostContent {
	T := i18n.MustTfunc(builder.siteLang())

	result := &PostContent{
		Model: post,

		Title: post.Title,
		Url:   node.Url,
	}

	year, _, day := post.CreatedAt.Date()
	result.Date = T("post_format_date", map[string]interface{}{
		"Year":  year,
		"Month": T("month_" + post.CreatedAt.Format("January")),
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
	node := builder.newNodeForKind(KIND_POSTS)
	node.fillUrl(slug)

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
