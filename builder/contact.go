package builder

// Builder for contact page
type ContactBuilder struct {
	*NodeBuilderBase
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
	node := builder.NewNode()

	node.Title = "Contact"
	node.Meta = &NodeMeta{Description: "Contact test node"}
	node.Content = "Soon"
	node.InNavBar = true

	builder.AddNode(node)
}
