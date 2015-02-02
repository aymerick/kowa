package builder

import (
	"html/template"
	"strings"
)

// Builder for contact page
type ContactBuilder struct {
	*NodeBuilderBase
}

// Contact content for template
type ContactContent struct {
	Title   string
	Tagline string

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
	title := "Contact"
	tagline := "" // @todo

	contactContent := builder.NewContactContent(title, tagline)
	if contactContent.HaveContact || contactContent.HaveSocial {
		node := builder.newNode()
		node.fillUrl(node.Kind)

		node.Title = title
		node.Meta = &NodeMeta{Description: tagline} // @todo !!!
		node.Content = contactContent
		node.InNavBar = true
		node.NavBarOrder = 10

		builder.addNode(node)
	}
}

func (builder *ContactBuilder) NewContactContent(title string, tagline string) *ContactContent {
	result := &ContactContent{
		Title:   title,
		Tagline: tagline,
	}

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
