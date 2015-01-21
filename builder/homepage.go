package builder

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

	node.Content = "Soon"

	builder.AddNode(node)
}
