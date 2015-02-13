package builder

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
)

// Homepage node builder
type HomepageBuilder struct {
	*NodeBuilderBase
}

// Homepage node content
type HomepageContent struct {
	Node *Node

	Description template.HTML // Site description
	MoreDesc    template.HTML // Site additional description
	JoinText    template.HTML // Site join text
	Logo        *ImageVars    // Site logo
	Cover       *ImageVars    // Site cover

	Activities []*ActivityVars // Activities
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
	node.Tagline = site.Tagline
	node.Meta = &NodeMeta{
		Title:       site.Name,
		Description: site.Tagline,
	}

	node.Content = builder.NewHomepageContent(node)

	builder.addNode(node)
}

// Instanciate a new homepage content
func (builder *HomepageBuilder) NewHomepageContent(node *Node) *HomepageContent {
	site := builder.site()

	result := &HomepageContent{
		Node: node,
	}

	sanitizePolicy := bluemonday.UGCPolicy()
	sanitizePolicy.AllowAttrs("style").OnElements("p", "span") // I know this is bad

	result.Description = template.HTML(sanitizePolicy.Sanitize(site.Description))
	result.MoreDesc = template.HTML(sanitizePolicy.Sanitize(site.MoreDesc))
	result.JoinText = template.HTML(sanitizePolicy.Sanitize(site.JoinText))

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
