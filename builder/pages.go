package builder

import (
	"html/template"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Builder for pages
type PagesBuilder struct {
	*NodeBuilderBase
}

// Page content for template
type PageContent struct {
	Title   string
	Tagline string

	Date  time.Time
	Cover string
	Body  template.HTML
	Url   string
}

func init() {
	RegisterNodeBuilder(KIND_PAGE, NewPagesBuilder)
}

func NewPagesBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &PagesBuilder{
		&NodeBuilderBase{
			nodeKind:    KIND_PAGE,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *PagesBuilder) Load() {
	for _, page := range *builder.site().FindAllPages() {
		builder.loadPage(page)
	}
}

// Build page
func (builder *PagesBuilder) loadPage(page *models.Page) {
	node := builder.newNode()
	node.fillUrl(utils.Urlify(page.Title))

	node.Title = page.Title
	node.Meta = &NodeMeta{Description: page.Tagline}
	node.Content = builder.NewPageContent(page, node)

	node.InNavBar = page.InNavBar

	builder.addNode(node)
}

// Instanciate a new page content
func (builder *PagesBuilder) NewPageContent(page *models.Page, node *Node) *PageContent {
	result := &PageContent{
		Title:   page.Title,
		Tagline: page.Tagline,
		Date:    page.CreatedAt,
		Url:     node.Url,
	}

	cover := page.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover, models.MEDIUM_KIND)
	}

	html := blackfriday.MarkdownCommon([]byte(page.Body))
	result.Body = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	return result
}
