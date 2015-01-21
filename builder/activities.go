package builder

// Builder for activities page
type ActivitiesBuilder struct {
	*NodeBuilder
}

func NewActivitiesBuilder(site *Site) *ActivitiesBuilder {
	return &ActivitiesBuilder{
		&NodeBuilder{
			NodeKind: KIND_ACTIVITIES,
			site:     site,
		},
	}
}

// NodeBuilderInterface
func (builder *ActivitiesBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Activities"

	node.Meta = &NodeMeta{
		Description: "Activities test page",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
