package builder

// Builder for pages
type PagesBuilder struct {
	*NodeBuilder
}

func NewPagesBuilder(site *Site) *PagesBuilder {
	return &PagesBuilder{
		&NodeBuilder{
			NodeKind: KIND_PAGE,
			site:     site,
		},
	}
}

// NodeBuilderInterface
func (builder *PagesBuilder) Load() {
	node := builder.NewNode()

	node.Path = "page-1.html"

	node.Title = "Page #1"

	node.Meta = &NodeMeta{
		Description: "Page test page #1",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
