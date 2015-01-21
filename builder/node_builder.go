package builder

import (
	"log"
	"os"
)

// node builder
type NodeBuilder struct {
	Nodes    []*Node
	NodeKind string

	site *Site
}

// interface for node builder
type NodeBuilderInterface interface {
	Site() *Site
	Load()
	Generate()
}

// NodeBuilderInterface
func (builder *NodeBuilder) Site() *Site {
	return builder.site
}

// NodeBuilderInterface
func (builder *NodeBuilder) Load() {
	panic("Should be implemented by includer")
}

// NodeBuilderInterface
func (builder *NodeBuilder) Generate() {
	for _, node := range builder.Nodes {
		builder.GenerateNode(node)
	}
}

// generate given node
func (builder *NodeBuilder) GenerateNode(node *Node) {
	osFilePath := builder.Site().FilePath(node.FullUrl())

	// ensure dir
	if err := builder.Site().EnsureFileDir(osFilePath); err != nil {
		builder.AddGenError(err)
		return
	}

	// open file
	outputFile, err := os.Create(osFilePath)
	if err != nil {
		builder.AddGenError(err)
		return
	}
	defer outputFile.Close()

	// write to file
	log.Printf("[DBG] Writing file: %s", osFilePath)
	if err := node.Generate(outputFile, builder.Site().Layout()); err != nil {
		builder.AddGenError(err)
	}
}

// init a new node with builder node kind
func (builder *NodeBuilder) NewNode() *Node {
	return builder.NewNodeForKind(builder.NodeKind)
}

// init a new node with given node kind
func (builder *NodeBuilder) NewNodeForKind(nodeKind string) *Node {
	return NewNode(builder, nodeKind)
}

// add a new node to build
func (builder *NodeBuilder) AddNode(node *Node) {
	builder.Nodes = append(builder.Nodes, node)
}

// add a node generation error
func (builder *NodeBuilder) AddGenError(err error) {
	builder.Site().AddGenError(builder.NodeKind, err)
}
