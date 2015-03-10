package builder

import (
	"html/template"
	"time"

	"github.com/aymerick/kowa/models"
)

// Page nodes builder
type PagesBuilder struct {
	*NodeBuilderBase
}

// Page node content
type PageContent struct {
	Node  *Node
	Model *models.Page

	Date  time.Time
	Cover *ImageVars
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
	node.fillUrl(page.Title)

	pageContent := builder.NewPageContent(page, node)
	if pageContent.Body != "" {
		node.Title = page.Title
		node.Tagline = page.Tagline
		node.Cover = pageContent.Cover

		node.Meta = &NodeMeta{Description: page.Tagline}
		node.InNavBar = page.InNavBar

		node.Content = pageContent

		builder.addNode(node)
	}
}

// Instanciate a new page content
func (builder *PagesBuilder) NewPageContent(page *models.Page, node *Node) *PageContent {
	result := &PageContent{
		Node:  node,
		Model: page,

		Date: page.CreatedAt,
		Url:  node.Url,
	}

	cover := page.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover)
	}

	result.Body = generateHTML(page.Format, page.Body)

	return result
}
