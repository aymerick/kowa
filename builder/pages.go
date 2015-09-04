package builder

import (
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/raymond"
)

// PagesBuilder builds custom pages
type PagesBuilder struct {
	*NodeBuilderBase
}

// PageContent represents a page node content
type PageContent struct {
	Model *models.Page

	Date  time.Time
	Cover *ImageVars
	Body  raymond.SafeString
	Url   string
}

func init() {
	RegisterNodeBuilder(kindPage, NewPagesBuilder)
}

// NewPagesBuilder instanciates a new NodeBuilder
func NewPagesBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &PagesBuilder{
		&NodeBuilderBase{
			nodeKind:    kindPage,
			siteBuilder: siteBuilder,
		},
	}
}

// Load is part of NodeBuilder interface
func (builder *PagesBuilder) Load() {
	for _, page := range *builder.site().FindAllPages() {
		builder.loadPage(page)
	}
}

// Build page
func (builder *PagesBuilder) loadPage(page *models.Page) {
	node := builder.newNode()
	node.fillURL(page.Title)

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

// NewPageContent instanciates a new PageContent
func (builder *PagesBuilder) NewPageContent(page *models.Page, node *Node) *PageContent {
	result := &PageContent{
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
