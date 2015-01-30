package builder

// Builder for contact page
type ContactBuilder struct {
	*NodeBuilderBase
}

// Contact content for template
type ContactContent struct {
	Title   string
	Tagline string
}

func init() {
	RegisterNodeBuilder(KIND_CONTACT, NewContactBuilder)
}

func NewContactBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &ContactBuilder{
		&NodeBuilderBase{
			nodeKind:    KIND_CONTACT,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *ContactBuilder) Load() {
	node := builder.newNode()
	node.fillUrl(node.Kind)

	title := "Contact"
	tagline := "" // @todo

	node.Title = title
	node.Meta = &NodeMeta{Description: tagline} // @todo !!!
	node.Content = &ContactContent{
		Title:   title,
		Tagline: tagline,
	}
	node.InNavBar = true

	builder.addNode(node)
}
