package builder

// Builder for contact page
type ContactBuilder struct {
	*NodeBuilderBase
}

func init() {
	RegisterBuilderInitializer(KIND_CONTACT, NewContactBuilder)
}

func NewContactBuilder(site *Site) NodeBuilder {
	return &ContactBuilder{
		&NodeBuilderBase{
			NodeKind: KIND_CONTACT,
			site:     site,
		},
	}
}

// NodeBuilder
func (builder *ContactBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Contact"

	node.Meta = &NodeMeta{
		Description: "Contact test node",
	}

	node.Content = "Soon"

	node.InNavBar = true

	builder.AddNode(node)
}
