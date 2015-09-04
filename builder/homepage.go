package builder

import (
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/raymond"
)

// HomepageBuilder builds homepage
type HomepageBuilder struct {
	*NodeBuilderBase
}

// HomepageContent represents homepage node content
type HomepageContent struct {
	Description raymond.SafeString // Site description
	MoreDesc    raymond.SafeString // Site additional description
	JoinText    raymond.SafeString // Site join text
	Logo        *ImageVars         // Site logo
	Cover       *ImageVars         // Site cover

	Activities []*ActivityVars // Activities
}

func init() {
	RegisterNodeBuilder(kindHomepage, NewHomepageBuilder)
}

// NewHomepageBuilder instanciates a new NodeBuilder
func NewHomepageBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &HomepageBuilder{
		&NodeBuilderBase{
			nodeKind:    kindHomepage,
			siteBuilder: siteBuilder,
		},
	}
}

// Load is part of NodeBuilder interface
func (builder *HomepageBuilder) Load() {
	T := i18n.MustTfunc(builder.siteLang())

	node := builder.newNode()
	node.fillURL("")

	site := builder.site()

	name := site.Name
	if name == "" {
		name = T("empty_site_name")
	}

	tagline := site.Tagline
	if tagline == "" {
		tagline = T("empty_site_tagline")
	}

	node.Title = name
	node.Tagline = tagline
	node.Meta = &NodeMeta{
		Title:       site.Name,
		Description: site.Tagline,
	}

	node.Content = builder.NewHomepageContent(node)

	builder.addNode(node)
}

// NewHomepageContent instanciates a new homepage content
func (builder *HomepageBuilder) NewHomepageContent(node *Node) *HomepageContent {
	site := builder.site()

	result := &HomepageContent{}

	result.Description = generateHTML(models.FORMAT_HTML, site.Description)
	result.MoreDesc = generateHTML(models.FORMAT_HTML, site.MoreDesc)
	result.JoinText = generateHTML(models.FORMAT_HTML, site.JoinText)

	logo := site.FindLogo()
	if logo != nil {
		result.Logo = builder.addImage(logo)
	}

	cover := site.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover)
	}

	result.Activities = builder.SiteBuilder().activitiesVars()
	if len(result.Activities) > 6 {
		result.Activities = result.Activities[0:6]
	}

	return result
}
