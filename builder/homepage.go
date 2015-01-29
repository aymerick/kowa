package builder

import (
	"html/template"

	"github.com/aymerick/kowa/models"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Builder for homepage
type HomepageBuilder struct {
	*NodeBuilderBase
}

// Homepage content for template
type HomepageContent struct {
	Name        string        // Site name
	Tagline     string        // Site tagline
	Description template.HTML // Site description
	MoreDesc    template.HTML // Site additional description
	JoinText    template.HTML // Site join text
	Logo        string        // Site logo
	Cover       string        // Site cover

	Activities []*ActivityContent // Activities
}

func init() {
	RegisterNodeBuilder(KIND_HOMEPAGE, NewHomepageBuilder)
}

func NewHomepageBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &HomepageBuilder{
		&NodeBuilderBase{
			nodeKind:    KIND_HOMEPAGE,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *HomepageBuilder) Load() {
	node := builder.newNode()
	node.fillUrl("")

	site := builder.site()

	node.Title = site.Name
	node.Meta = &NodeMeta{
		Title:       site.Name,
		Description: site.Tagline,
	}
	node.Content = builder.NewHomepageContent()

	builder.addNode(node)
}

// Instanciate a new homepage content
func (builder *HomepageBuilder) NewHomepageContent() *HomepageContent {
	site := builder.site()

	result := &HomepageContent{
		Name:    site.Name,
		Tagline: site.Tagline,
	}

	var html []byte

	html = blackfriday.MarkdownCommon([]byte(site.Description))
	result.Description = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	html = blackfriday.MarkdownCommon([]byte(site.MoreDesc))
	result.MoreDesc = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	html = blackfriday.MarkdownCommon([]byte(site.JoinText))
	result.JoinText = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	logo := site.FindLogo()
	if logo != nil {
		result.Logo = builder.addImage(logo, models.MEDIUM_KIND)
	}

	cover := site.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover, models.MEDIUM_KIND)
	}

	result.Activities = builder.SiteBuilder().activitiesContents()
	if len(result.Activities) > 6 {
		result.Activities = result.Activities[0:6]
	}

	return result
}
