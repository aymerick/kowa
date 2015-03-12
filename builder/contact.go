package builder

import (
	"html/template"
	"strings"

	"github.com/aymerick/kowa/utils"
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
	Address     template.HTML

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
	T := i18n.MustTfunc(utils.DEFAULT_LANG) // @todo i18n

	title := T("contact")
	tagline := "" // @todo Fill

	contactContent := builder.NewContactContent()
	if contactContent.HaveContact || contactContent.HaveSocial {
		node := builder.newNode()
		node.fillUrl(title)

		node.Title = title
		node.Tagline = tagline
		node.Meta = &NodeMeta{Description: tagline} // @todo !!!
		node.InNavBar = true
		node.NavBarOrder = 20

		contactContent.Node = node
		node.Content = contactContent

		builder.addNode(node)
	}
}

func (builder *ContactBuilder) NewContactContent() *ContactContent {
	result := &ContactContent{}

	site := builder.site()

	result.Email = site.Email

	addrSafe := template.HTMLEscapeString(site.Address)
	result.Address = template.HTML(strings.Replace(addrSafe, "\n", "<br />\n", -1))

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
