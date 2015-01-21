package builder

import (
	"errors"
	"fmt"
	"os"
	"path"
)

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
func (builder *NodeBuilder) Generate() {
	for _, node := range builder.Nodes {
		builder.GenerateNode(node)
	}
}

func (builder *NodeBuilder) GenerateNode(node *Node) {
	filePath := node.FilePath()
	if filePath == "" {
		builder.AddGenError(errors.New(fmt.Sprintf("No path defined for node: %v", node)))
		return
	}

	outputFile, err := os.Create(path.Join(builder.Site.GenDir, filePath))
	if err != nil {
		builder.AddGenError(err)
		return
	}

	if err := node.Generate(outputFile, builder.Site.Layout); err != nil {
		builder.AddGenError(err)
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
func (builder *NodeBuilder) AddGenError(err error) {
	builder.Site.AddGenError(builder.NodeKind, err)
}