package builder

import (
	"errors"
	"html/template"
	"io"
	"io/ioutil"
)

// site page metadata
type SitePageMeta struct {
	Description string
}

// site page
type SitePage struct {
	Kind      string
	Title     string
	Meta      *SitePageMeta
	BodyClass string
	Head      string
	Footer    string
	Content   interface{}

	layout   *template.Template
	template *template.Template
}

// interface for site page builders
type SitePageBuilder interface {
	Fill(page *SitePage, site *Site) error
}

// all site page kinds
const (
	KIND_ACTIVITIES = "activities"
	KIND_CONTACT    = "contact"
	KIND_EVENT      = "event"
	KIND_EVENTS     = "events"
	KIND_INDEX      = "index"
	KIND_MEMBERS    = "members"
	KIND_PAGE       = "page"
	KIND_POST       = "post"
	KIND_POSTS      = "posts"
)

// site page builders
var SitePageBuilders = map[string]SitePageBuilder{
	KIND_ACTIVITIES: NewActivitiesBuilder(),
	KIND_CONTACT:    nil,
	KIND_EVENT:      nil,
	KIND_EVENTS:     nil,
	KIND_INDEX:      nil,
	KIND_MEMBERS:    nil,
	KIND_PAGE:       nil,
	KIND_POST:       nil,
	KIND_POSTS:      nil,
}

// create a new site page
func NewSitePage(kind string) *SitePage {
	return &SitePage{
		Kind: kind,
	}
}

// get site page template
func (page *SitePage) Template() (*template.Template, error) {
	if page.template != nil {
		return page.template, nil
	} else if page.layout == nil {
		return nil, errors.New("Can't generate page without a layout template")
	} else {
		result := template.Must(page.layout.Clone())

		binData, err := ioutil.ReadFile(templatePath(page.Kind))
		if err == nil {
			_, err = result.New("content").Parse(string(binData))
			if err != nil {
				page.template = result
			}
		}

		return result, err
	}
}

// generate site page
func (page *SitePage) Generate(wr io.Writer) error {
	tpl := template.Must(page.Template())

	return tpl.Execute(wr, page)
}
