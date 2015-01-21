package builder

type ContactBuilder struct {
	*NodeBuilder
}

func NewContactBuilder(site *Site) *ContactBuilder {
	return &ContactBuilder{
		&NodeBuilder{
			Site:     site,
			NodeKind: KIND_CONTACT,
		},
	}
}

func (builder *ContactBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Contact"

	node.Meta = &NodeMeta{
		Description: "Contact test node",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
