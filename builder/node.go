package builder

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
)

// Node
type Node struct {
	// template vars
	Kind string
	Site *SiteVars

	Title     string
	Meta      *NodeMeta
	BodyClass string
	Content   interface{}
	InNavBar  bool
	Slug      string
	Url       string
	FullUrl   string

	builder  NodeBuilder
	template *template.Template
}

// Node metadata
type NodeMeta struct {
	Description string
}

// All node kinds
const (
	KIND_ACTIVITIES = "activities"
	KIND_CONTACT    = "contact"
	KIND_HOMEPAGE   = "homepage"
	KIND_PAGE       = "page"
	KIND_POST       = "post"
	KIND_POSTS      = "posts"
)

// Create a new node
func NewNode(builder NodeBuilder, kind string) *Node {
	return &Node{
		Kind:      kind,
		BodyClass: kind,
		InNavBar:  false,

		builder: builder,
	}
}

func (node *Node) fillUrl(slug string) {
	// Slug
	node.Slug = slug

	// FullUrl
	if node.builder.SiteBuilder().uglyURL || (node.Slug == "") || (node.Slug == "/") || (node.Slug == "index") {
		name := node.Slug
		switch name {
		case "", "/":
			name = "index"
		}

		// ugly URL (or homepage)
		node.FullUrl = path.Join("/", fmt.Sprintf("%s.html", name))
	} else {
		// pretty URL
		node.FullUrl = path.Join("/", node.Slug, "index.html")
	}

	// Url
	dir, fileName := path.Split(node.FullUrl)
	if fileName == "index.html" {
		// pretty URL
		node.Url = dir
	} else {
		// ugly URL
		node.Url = node.FullUrl
	}
}

// Get node template
func (node *Node) Template(layout *template.Template) (*template.Template, error) {
	if node.template != nil {
		return node.template, nil
	} else if layout == nil {
		return nil, errors.New("Can't generate node without a layout template")
	} else {
		result := template.Must(layout.Clone())

		binData, err := ioutil.ReadFile(node.builder.SiteBuilder().templatePath(node.Kind))
		if err == nil {
			_, err = result.New("content").Parse(string(binData))
			if err != nil {
				node.template = result
			}
		}

		return result, err
	}
}

// Generate node
func (node *Node) Generate(wr io.Writer, layout *template.Template) error {
	tpl := template.Must(node.Template(layout))

	return tpl.Execute(wr, node)
}
