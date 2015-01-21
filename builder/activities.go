package builder

type ActivitiesBuilder struct {
	*NodeBuilder
}

func NewActivitiesBuilder(site *Site) *ActivitiesBuilder {
	return &ActivitiesBuilder{
		&NodeBuilder{
			Site:     site,
			NodeKind: KIND_ACTIVITIES,
		},
	}
}

func (builder *ActivitiesBuilder) Load() {
	node := builder.NewNode()

	node.Title = "Activities"

	node.Meta = &NodeMeta{
		Description: "Activities test page",
	}

	node.Content = "Soon"

	builder.AddNode(node)
}
