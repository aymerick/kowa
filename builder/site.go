package builder

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
)

type Site struct {
	Layout *template.Template
	Pages  map[string]*SitePage
	Errors []error
}

func NewSite() *Site {
	result := &Site{
		Layout: template.Must(template.ParseFiles(templatePath("layout"))),
		Pages:  make(map[string]*SitePage),
		Errors: []error{},
	}

	// @todo Load partials
	// template.Must(site.Layout.ParseGlob(path.Join(partialsPath(), "*.html")))

	return result
}

// Build site
func (site *Site) Build() {
	// load pages
	if len(site.Pages) == 0 {
		site.LoadPages()
		if len(site.Errors) > 0 {
			log.Printf("%d error(s) while loading pages: %v", len(site.Errors), site.Errors)
			site.Errors = nil
		}
	}

	// fill pages data
	site.FillPages()
	if len(site.Errors) > 0 {
		log.Printf("%d error(s) while filling pages: %v", len(site.Errors), site.Errors)
		site.Errors = nil
	}

	// generate pages
	site.GeneratePages()
	if len(site.Errors) > 0 {
		log.Printf("%d error(s) while generating pages: %v", len(site.Errors), site.Errors)
		site.Errors = nil
	}
}

// Load site pages
func (site *Site) LoadPages() {
	for kind, _ := range SitePageBuilders {
		site.Pages[kind] = NewSitePage(kind)
		site.Pages[kind].layout = site.Layout
	}
}

// Fill site pages
func (site *Site) FillPages() {
	for kind, page := range site.Pages {
		if SitePageBuilders[kind] != nil {
			if err := SitePageBuilders[kind].Fill(page, site); err != nil {
				site.Errors = append(site.Errors, err)
			}
		}
	}
}

// Generate site pages
func (site *Site) GeneratePages() {
	for kind, page := range site.Pages {
		if SitePageBuilders[kind] == nil {
			site.Errors = append(site.Errors, errors.New(fmt.Sprintf("Can't generate page %s because there is no builder for it", kind)))
		} else if err := page.Generate(os.Stdout); err != nil {
			site.Errors = append(site.Errors, err)
		}
	}
}
