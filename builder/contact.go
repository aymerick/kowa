package builder

import (
	"strings"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/raymond"
	"github.com/nicksnyder/go-i18n/i18n"
)

// Contact node builder
type ContactBuilder struct {
	*NodeBuilderBase
}

// Contact node content
type ContactContent struct {
	Node *Node

	HaveContact bool
	Email       string
	Address     raymond.SafeString

	HaveSocial bool
	Facebook   string
	Twitter    string
	GooglePlus string
}

func init() {
	RegisterNodeBuilder(KIND_CONTACT, NewContactBuilder)
}

func NewContactBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &ContactBuilder{
		&NodeBuilderBase{
			nodeKind:    KIND_CONTACT,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *ContactBuilder) Load() {
	// get node content
	contactContent := builder.NewContactContent()
	if !contactContent.HaveContact && !contactContent.HaveSocial {
		return
	}

	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PAGE_KIND_CONTACT)
	if disabled {
		return
	}

	T := i18n.MustTfunc(builder.siteLang())
	slug := T("contact")

	if title == "" {
		title = slug
	}

	// build node
	node := builder.newNode()
	node.fillUrl(slug)

	node.Title = title
	node.Tagline = tagline
	node.Cover = cover

	node.Meta = &NodeMeta{Description: tagline}

	node.InNavBar = true
	node.NavBarOrder = 20

	contactContent.Node = node
	node.Content = contactContent

	builder.addNode(node)
}

func (builder *ContactBuilder) NewContactContent() *ContactContent {
	result := &ContactContent{}

	site := builder.site()

	result.Email = site.Email

	addrSafe := raymond.Escape(site.Address)
	result.Address = raymond.SafeString(strings.Replace(addrSafe, "\n", "<br />\n", -1))

	if result.Email != "" || result.Address != "" {
		result.HaveContact = true
	}

	result.Facebook = site.Facebook
	result.Twitter = site.Twitter
	result.GooglePlus = site.GooglePlus

	if result.Facebook != "" || result.Twitter != "" || result.GooglePlus != "" {
		result.HaveSocial = true
	}

	return result
}
