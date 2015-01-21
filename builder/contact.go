package builder

// Builder for contact page
type ContactBuilder struct {
	*NodeBuilder
}

func NewContactBuilder(site *Site) *ContactBuilder {
	return &ContactBuilder{
		&NodeBuilder{
			NodeKind: KIND_CONTACT,
			site:     site,
		},
	}
}

// NodeBuilderInterface
func (builder *ContactBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Contact"

	node.Meta = &NodeMeta{
		Description: "Contact test node",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
