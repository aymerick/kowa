package builder

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
)

type Site struct {
	Layout    *template.Template
	SitePages map[string]*SitePage
	Errors    []error
}

func NewSite() *Site {
	result := &Site{
		Layout:    template.Must(template.ParseFiles(templatePath("layout"))),
		SitePages: make(map[string]*SitePage),
		Errors:    []error{},
	}

	// @todo Load partials
	// template.Must(site.Layout.ParseGlob(path.Join(partialsPath(), "*.html")))

	return result
}

// Build site
func (site *Site) Build() {
	// load
	if len(site.SitePages) == 0 {
		site.LoadSitePages()
		if len(site.Errors) > 0 {
			log.Printf("%d error(s) while loading site pages: %v", len(site.Errors), site.Errors)
			site.Errors = nil
		}
	}

	// fill
	site.FillSitePages()
	if len(site.Errors) > 0 {
		log.Printf("%d error(s) while filling site pages: %v", len(site.Errors), site.Errors)
		site.Errors = nil
	}

	// generate
	site.GenerateSitePages()
	if len(site.Errors) > 0 {
		log.Printf("%d error(s) while generating site pages: %v", len(site.Errors), site.Errors)
		site.Errors = nil
	}
}

// Load site site pages
func (site *Site) LoadSitePages() {
	for kind, _ := range SitePageBuilders {
		site.SitePages[kind] = NewSitePage(kind)
		site.SitePages[kind].layout = site.Layout
	}
}

// Fill site pages
func (site *Site) FillSitePages() {
	for kind, sitePage := range site.SitePages {
		sitePage.BodyClass = kind

		if SitePageBuilders[kind] != nil {
			if err := SitePageBuilders[kind].Fill(sitePage, site); err != nil {
				site.Errors = append(site.Errors, err)
			}
		}
	}
}

// Generate site pages
func (site *Site) GenerateSitePages() {
	for kind, sitePage := range site.SitePages {
		if SitePageBuilders[kind] == nil {
			site.Errors = append(site.Errors, errors.New(fmt.Sprintf("Can't generate site page %s because there is no builder for it", kind)))
		} else {
			// generate
			if err := sitePage.Generate(os.Stdout); err != nil {
				site.Errors = append(site.Errors, err)
			}
		}
	}
}
