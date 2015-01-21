package builder

type PagesBuilder struct {
	*NodeBuilder
}

func NewPagesBuilder(site *Site) *PagesBuilder {
	return &PagesBuilder{
		&NodeBuilder{
			Site:     site,
			NodeKind: KIND_PAGE,
		},
	}
}

func (builder *PagesBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Page #1"

	node.Meta = &NodeMeta{
		Description: "Page test page #1",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
