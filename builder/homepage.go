package builder

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Homepage content for template
type HomepageContent struct {
	Name        string        // Site name
	Tagline     string        // Site tagline
	Description template.HTML // Site description
	MoreDesc    template.HTML // Site additional description
	JoinText    template.HTML // Site join text
	Logo        string        // Site logo
	Cover       string        // Site cover
}

// Builder for homepage
type HomepageBuilder struct {
	*NodeBuilder
}

func NewHomepageBuilder(site *Site) *HomepageBuilder {
	return &HomepageBuilder{
		&NodeBuilder{
			NodeKind: KIND_HOMEPAGE,
			site:     site,
		},
	}
}

// NodeBuilderInterface
func (builder *HomepageBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Homepage"

	node.Meta = &NodeMeta{
		Description: "Homepage test node",
	}

	node.Content = builder.NewHomepageContent()

	builder.AddNode(node)
}

/// Instanciate a new homepage content
func (builder *HomepageBuilder) NewHomepageContent() *HomepageContent {
	siteModel := builder.Site().Model

	result := &HomepageContent{
		Name:    siteModel.Name,
		Tagline: siteModel.Tagline,
	}

	var html []byte

	html = blackfriday.MarkdownCommon([]byte(siteModel.Description))
	result.Description = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	html = blackfriday.MarkdownCommon([]byte(siteModel.MoreDesc))
	result.MoreDesc = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	html = blackfriday.MarkdownCommon([]byte(siteModel.JoinText))
	result.JoinText = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	logo := siteModel.FindLogo()
	if logo != nil {
		result.Logo = logo.MediumURL()
	}

	cover := siteModel.FindCover()
	if cover != nil {
		result.Cover = cover.MediumURL()
	}

	return result
}
