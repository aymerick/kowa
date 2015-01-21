package builder

import (
	"html/template"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Page content for template
type PageContent struct {
	Date    time.Time     // CreatedAt
	Image   string        // Cover
	Title   string        // Title
	Tagline string        // Tagline
	Content template.HTML // Body
	Url     string        // Absolute URL
}

// Builder for pages
type PagesBuilder struct {
	*NodeBuilder
}

func NewPagesBuilder(site *Site) *PagesBuilder {
	return &PagesBuilder{
		&NodeBuilder{
			NodeKind: KIND_PAGE,
			site:     site,
		},
	}
}

// NodeBuilderInterface
func (builder *PagesBuilder) Load() {
	for _, page := range *builder.Site().Model.FindAllPages() {
		builder.BuildPage(page)
	}
}

// Build page
func (builder *PagesBuilder) BuildPage(page *models.Page) {
	node := builder.NewNode()

	node.basePath = utils.Urlify(page.Title)

	node.Title = page.Title
	node.Meta = &NodeMeta{
		Description: page.Tagline,
	}

	node.Content = builder.NewPageContent(page, node)

	builder.AddNode(node)
}

/// Instanciate a new page content
func (builder *PagesBuilder) NewPageContent(page *models.Page, node *Node) *PageContent {
	result := &PageContent{
		Date:    page.CreatedAt,
		Title:   page.Title,
		Tagline: page.Tagline,
		Url:     node.Url(),
	}

	cover := page.FindCover()
	if cover != nil {
		result.Image = cover.MediumURL()
	}

	html := blackfriday.MarkdownCommon([]byte(page.Body))
	result.Content = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	return result
}
