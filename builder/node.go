package builder

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
)

// node metadata
type NodeMeta struct {
	Description string
}

// node
type Node struct {
	Kind string

	Title     string
	Meta      *NodeMeta
	BodyClass string
	Head      string
	Footer    string
	Content   interface{}

	Builder NodeBuilderInterface

	slug     string
	template *template.Template
}

// all node kinds
const (
	KIND_ACTIVITIES = "activities"
	KIND_CONTACT    = "contact"
	KIND_HOMEPAGE   = "homepage"
	KIND_PAGE       = "page"
	KIND_POST       = "post"
	KIND_POSTS      = "posts"
)

// create a new node
func NewNode(builder NodeBuilderInterface, kind string) *Node {
	return &Node{
		Kind:    kind,
		Builder: builder,
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

	if node.Builder.Site().UglyURL || (slug == "index") {
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

// get node template
func (node *Node) Template(layout *template.Template) (*template.Template, error) {
	if node.template != nil {
		return node.template, nil
	} else if layout == nil {
		return nil, errors.New("Can't generate node without a layout template")
	} else {
		result := template.Must(layout.Clone())

		binData, err := ioutil.ReadFile(node.Builder.Site().TemplatePath(node.Kind))
		if err == nil {
			_, err = result.New("content").Parse(string(binData))
			if err != nil {
				node.template = result
			}
		}

		return result, err
	}
}

// generate node
func (node *Node) Generate(wr io.Writer, layout *template.Template) error {
	tpl := template.Must(node.Template(layout))

	return tpl.Execute(wr, node)
}
