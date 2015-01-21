package builder

import (
	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
)

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

func (builder *PagesBuilder) BuildPage(page *models.Page) {
	node := builder.NewNode()

	node.basePath = utils.Urlify(page.Title)

	node.Title = page.Title

	node.Meta = &NodeMeta{
		Description: page.Tagline,
	}

	node.Content = page

	builder.AddNode(node)
}
