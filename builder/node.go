package builder

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"

	"github.com/aymerick/kowa/utils"
)

// Node vars
type Node struct {
	// template vars
	Kind string
	Site *SiteVars

	Title   string
	Tagline string
	Cover   *ImageVars

	Meta        *NodeMeta
	BodyClass   string
	InNavBar    bool
	NavBarOrder int

	Slug        string // eg: 2015/03/17/my_post
	FilePath    string // eg: /2015/03/17/my_post/index.html
	Url         string // eg: /my_site/2015/03/17/my_post/
	AbsoluteUrl string // eg: http://127.0.0.1:48910/my_site/2015/03/17/my_post/

	Content interface{}

	builder NodeBuilder
}

// Node metadata
type NodeMeta struct {
	Title       string
	Description string
	ImageUrl    string

	// @todo article:published_time and article:modified_time

	OGType      string
	TwitterCard string
}

// All node kinds
const (
	KIND_ACTIVITIES = "activities"
	KIND_MEMBERS    = "members"
	KIND_CONTACT    = "contact"
	KIND_HOMEPAGE   = "homepage"
	KIND_PAGE       = "page"
	KIND_POST       = "post"
	KIND_POSTS      = "posts"
	KIND_EVENT      = "event"
	KIND_EVENTS     = "events"
)

// Create a new node
func NewNode(builder NodeBuilder, kind string) *Node {
	return &Node{
		Kind:        kind,
		BodyClass:   kind,
		NavBarOrder: 100,

		builder: builder,
	}
}

// Fill node Url
func (node *Node) fillUrl(slug string) {
	siteBuilder := node.builder.SiteBuilder()

	// Slug
	node.Slug = siteBuilder.addNodeSlug(utils.Pathify(slug))

	// FilePath
	if siteBuilder.site.UglyUrl || (node.Slug == "") || (node.Slug == "/") || (node.Slug == "index") {
		name := node.Slug
		switch name {
		case "", "/":
			name = "index"
		}

		// ugly URL (or homepage)
		node.FilePath = path.Join("/", fmt.Sprintf("%s.html", name))
	} else {
		// pretty URL
		node.FilePath = path.Join("/", node.Slug, "index.html")
	}

	var lastPart string

	dir, fileName := path.Split(node.FilePath)
	if fileName == "index.html" {
		// pretty URL
		lastPart = dir
	} else {
		// ugly URL
		lastPart = node.FilePath
	}

	// Url
	node.Url = utils.Urlify(path.Join(siteBuilder.basePath(), lastPart))

	// AbsoluteUrl
	node.AbsoluteUrl = utils.Urlify(fmt.Sprintf("%s%s", siteBuilder.site.BaseUrl, lastPart))
}

// Compute node template
func (node *Node) template(layout *template.Template) (*template.Template, error) {
	if layout == nil {
		return nil, errors.New("Can't generate node without a layout template")
	} else {
		result := template.Must(layout.Clone())

		// add "content" template to main layout
		binData, err := ioutil.ReadFile(node.builder.SiteBuilder().templatePath(node.Kind))
		if err == nil {
			_, err = result.New("content").Parse(string(binData))
		}

		return result, err
	}
}

// Generate node
func (node *Node) generate(wr io.Writer, layout *template.Template) error {
	tpl := template.Must(node.template(layout))

	return tpl.Execute(wr, node)
}
