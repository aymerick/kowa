package builder

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
)

// Node metadata
type NodeMeta struct {
	Description string
}

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

	builder  NodeBuilder
	slug     string
	template *template.Template
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

func (node *Node) Slug() string {
	if node.slug == "" {
		if node.Kind == KIND_HOMEPAGE {
			return "index"
		} else {
			return node.Kind
		}
	} else {
		return node.slug
	}
}

func (node *Node) FullUrl() string {
	slug := node.Slug()

	if node.builder.SiteBuilder().uglyURL || (slug == "index") {
		return path.Join("/", fmt.Sprintf("%s.html", slug))
	} else {
		return path.Join("/", slug, "index.html")
	}
}

func (node *Node) Url() string {
	fullUrl := node.FullUrl()

	dir, fileName := path.Split(fullUrl)
	if fileName == "index.html" {
		return dir
	} else {
		return fullUrl
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
