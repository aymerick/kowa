package builder

import (
	"fmt"
	"os"

	"github.com/aymerick/kowa/models"
)

// Interface for node builder
type NodeBuilder interface {
	// Returns site builder
	SiteBuilder() *SiteBuilder

	// Load all nodes
	Load()

	// Generate all nodes
	Generate()

	// Returns all loaded nodes
	Nodes() []*Node

	// Returns loaded nodes that must be placed in navigation bar
	NavBarNodes() []*Node

	// Returns given data
	Data(string) interface{}
}

// Node builder base
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
	panic("Must be implemented by includer")
}

// NodeBuilder
func (builder *NodeBuilderBase) Generate() {
	for _, node := range builder.nodes {
		// fill node with more data
		builder.fillNodeBeforeGeneration(node)

		builder.generateNode(node)
	}
}

// NodeBuilder
func (builder *NodeBuilderBase) Nodes() []*Node {
	return builder.nodes
}

// NodeBuilder
func (builder *NodeBuilderBase) NavBarNodes() []*Node {
	return builder.navBarNodes
}

// NodeBuilder
func (builder *NodeBuilderBase) Data(name string) interface{} {
	// Should be implemented by includer
	return nil
}

// Get site model
func (builder *NodeBuilderBase) site() *models.Site {
	return builder.SiteBuilder().site
}

// Fill node with more data
func (builder *NodeBuilderBase) fillNodeBeforeGeneration(node *Node) {
	node.Site = builder.siteBuilder.siteVars

	// defaults
	if node.Meta.Title == "" {
		// @todo Filter characters ?
		node.Meta.Title = fmt.Sprintf("%s - %s", node.Title, builder.site().Name)
	}
}

// Generate given node
func (builder *NodeBuilderBase) generateNode(node *Node) {
	osFilePath := builder.siteBuilder.filePath(node.FilePath)

	// ensure dir
	if err := builder.siteBuilder.ensureFileDir(osFilePath); err != nil {
		builder.addError(err)
		return
	}

	// open file
	outputFile, err := os.Create(osFilePath)
	if err != nil {
		builder.addError(err)
		return
	}
	defer outputFile.Close()

	// write to file
	// log.Printf("[DBG] Writing file: %s", osFilePath)
	if err := node.generate(outputFile, builder.siteBuilder.layout()); err != nil {
		builder.addError(err)
	}
}

// Init a new node with builder node kind
func (builder *NodeBuilderBase) newNode() *Node {
	return builder.newNodeForKind(builder.nodeKind)
}

// Init a new node with given node kind
func (builder *NodeBuilderBase) newNodeForKind(nodeKind string) *Node {
	return NewNode(builder, nodeKind)
}

// Add a new node to build
func (builder *NodeBuilderBase) addNode(node *Node) {
	builder.nodes = append(builder.nodes, node)

	if node.InNavBar {
		builder.navBarNodes = append(builder.navBarNodes, node)
	}
}

// Add an image to copy
func (builder *NodeBuilderBase) addImage(img *models.Image) *ImageVars {
	return builder.siteBuilder.addImage(img)
}

// Add a node generation error
func (builder *NodeBuilderBase) addError(err error) {
	builder.siteBuilder.addNodeBuilderError(builder.nodeKind, err)
}
