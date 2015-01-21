package builder

import "io"

// node builder
type NodeBuilder struct {
	Site     *Site
	Nodes    []*Node
	NodeKind string
}

// interface for node builder
type NodeBuilderInterface interface {
	Load()
	Generate()
}

// generate nodes
func (builder *NodeBuilder) Generate(wr io.Writer) {
	for _, node := range builder.Nodes {
		if err := node.Generate(wr, builder.Site.Layout); err != nil {
			builder.AddGenerationError(err)
		}
	}
}

// init a new node with builder node kind
func (builder *NodeBuilder) NewNode() *Node {
	return builder.NewNodeForKind(builder.NodeKind)
}

// init a new node with given node kind
func (builder *NodeBuilder) NewNodeForKind(nodeKind string) *Node {
	return NewNode(nodeKind)
}

// add a new node to build
func (builder *NodeBuilder) AddNode(node *Node) {
	builder.Nodes = append(builder.Nodes, node)
}

// add a node generation error
func (builder *NodeBuilder) AddGenerationError(err error) {
	builder.Site.AddGenerationError(builder.NodeKind, err)
}
