package builder

import (
	"os"

	"github.com/aymerick/kowa/models"
)

// interface for node builder
type NodeBuilder interface {
	// Returns site builder
	SiteBuilder() *SiteBuilder

	// Load all nodes
	Load()

	// Generate all nodes
	Generate()

	// Returns loaded nodes that must be placed in navigation bar
	NavBarNodes() []*Node
}

// node builder base
type NodeBuilderBase struct {
	// All loaded nodes
	nodes []*Node

	// Loaded nodes that must be placed in navigation bar
	navBarNodes []*Node

	// Node kind
	nodeKind string

	// Suite builder
	siteBuilder *SiteBuilder
}

// NodeBuilder
func (builder *NodeBuilderBase) SiteBuilder() *SiteBuilder {
	return builder.siteBuilder
}

// NodeBuilder
func (builder *NodeBuilderBase) Load() {
	panic("Should be implemented by includer")
}

// NodeBuilder
func (builder *NodeBuilderBase) Generate() {
	for _, node := range builder.nodes {
		// fill node with more data
		node.Site = builder.siteBuilder.siteVars

		builder.generateNode(node)
	}
}

// NodeBuilder
func (builder *NodeBuilderBase) NavBarNodes() []*Node {
	return builder.navBarNodes
}

// generate given node
func (builder *NodeBuilderBase) generateNode(node *Node) {
	osFilePath := builder.siteBuilder.filePath(node.FullUrl())

	// ensure dir
	if err := builder.siteBuilder.ensureFileDir(osFilePath); err != nil {
		builder.addGenError(err)
		return
	}

	// open file
	outputFile, err := os.Create(osFilePath)
	if err != nil {
		builder.addGenError(err)
		return
	}
	defer outputFile.Close()

	// write to file
	// log.Printf("[DBG] Writing file: %s", osFilePath)
	if err := node.Generate(outputFile, builder.siteBuilder.layout()); err != nil {
		builder.addGenError(err)
	}
}

// init a new node with builder node kind
func (builder *NodeBuilderBase) NewNode() *Node {
	return builder.NewNodeForKind(builder.nodeKind)
}

// init a new node with given node kind
func (builder *NodeBuilderBase) NewNodeForKind(nodeKind string) *Node {
	return NewNode(builder, nodeKind)
}

// add a new node to build
func (builder *NodeBuilderBase) AddNode(node *Node) {
	builder.nodes = append(builder.nodes, node)

	if node.InNavBar {
		builder.navBarNodes = append(builder.navBarNodes, node)
	}
}

// add an image to copy
func (builder *NodeBuilderBase) AddImage(img *models.Image, kind string) string {
	return builder.siteBuilder.AddImage(img, kind)
}

// add a node generation error
func (builder *NodeBuilderBase) addGenError(err error) {
	builder.siteBuilder.addGenError(builder.nodeKind, err)
}
