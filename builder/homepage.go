package builder

type HomepageBuilder struct {
	*NodeBuilder
}

func NewHomepageBuilder(site *Site) *HomepageBuilder {
	return &HomepageBuilder{
		&NodeBuilder{
			Site:     site,
			NodeKind: KIND_HOMEPAGE,
		},
	}
}

func (builder *HomepageBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Homepage"

	node.Meta = &NodeMeta{
		Description: "Homepage test node",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
