package builder

import (
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/raymond"
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

	Type        string
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
	node.Slug = siteBuilder.addNodeSlug(helpers.Pathify(slug))

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
	node.Url = helpers.Urlify(path.Join(siteBuilder.basePath(), lastPart))

	// AbsoluteUrl
	node.AbsoluteUrl = helpers.Urlify(fmt.Sprintf("%s%s", siteBuilder.site.BaseUrl, lastPart))
}

// Compute node template
func (node *Node) template(layout *raymond.Template) (*raymond.Template, error) {
	if layout == nil {
		return nil, errors.New("Can't generate node without a layout template")
	}

	result := layout.Clone()

	filePath := node.builder.SiteBuilder().templatePath(node.Kind)

	if err := result.RegisterPartialFile(filePath, "content"); err != nil {
		return nil, err
	}

	return result, nil
}

// Generate node
func (node *Node) generate(wr io.Writer, layout *raymond.Template) error {
	tpl, err := node.template(layout)
	if err != nil {
		return err
	}

	output, err := tpl.Exec(node)
	if err != nil {
		return err
	}

	if _, err := wr.Write([]byte(output)); err != nil {
		return err
	}

	return nil
}
