package builder

// Builder for activities page
type ActivitiesBuilder struct {
	*NodeBuilderBase
}

func init() {
	RegisterNodeBuilder(KIND_ACTIVITIES, NewActivitiesBuilder)
}

func NewActivitiesBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &ActivitiesBuilder{
		&NodeBuilderBase{
			nodeKind:    KIND_ACTIVITIES,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *ActivitiesBuilder) Load() {
	node := builder.newNode()
	node.fillUrl(node.Kind)

	node.Title = "Activities"
	node.Meta = &NodeMeta{Description: "Activities test page"}
	node.Content = "Soon"
	node.InNavBar = true

	builder.addNode(node)
}
