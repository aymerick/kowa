package builder

import (
	"os"

	"github.com/aymerick/kowa/models"
)

// interface for node builder
type NodeBuilder interface {
	Site() *Site
	Load()
	Generate()
}

// node builder base
type NodeBuilderBase struct {
	Nodes    []*Node
	NodeKind string

	site   *Site
	images map[string]*models.Image
}

// NodeBuilder
func (builder *NodeBuilderBase) Site() *Site {
	return builder.site
}

// NodeBuilder
func (builder *NodeBuilderBase) Load() {
	panic("Should be implemented by includer")
}

// NodeBuilder
func (builder *NodeBuilderBase) Generate() {
	for _, node := range builder.Nodes {
		builder.GenerateNode(node)
	}
}

// generate given node
func (builder *NodeBuilderBase) GenerateNode(node *Node) {
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
	// log.Printf("[DBG] Writing file: %s", osFilePath)
	if err := node.Generate(outputFile, builder.Site().Layout()); err != nil {
		builder.AddGenError(err)
	}
}

// init a new node with builder node kind
func (builder *NodeBuilderBase) NewNode() *Node {
	return builder.NewNodeForKind(builder.NodeKind)
}

// init a new node with given node kind
func (builder *NodeBuilderBase) NewNodeForKind(nodeKind string) *Node {
	return NewNode(builder, nodeKind)
}

// add a new node to build
func (builder *NodeBuilderBase) AddNode(node *Node) {
	builder.Nodes = append(builder.Nodes, node)
}

// add an image to copy
func (builder *NodeBuilderBase) AddImage(img *models.Image, kind string) string {
	return builder.Site().AddImage(img, kind)
}

// add a node generation error
func (builder *NodeBuilderBase) AddGenError(err error) {
	builder.Site().AddGenError(builder.NodeKind, err)
}
