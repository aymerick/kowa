package builder

import (
	"fmt"
	"html/template"
	"os"
)

type Site struct {
	Layout         *template.Template
	ErrorCollector *ErrorCollector
}

func NewSite() *Site {
	result := &Site{
		Layout:         template.Must(template.ParseFiles(templatePath("layout"))),
		ErrorCollector: NewErrorCollector(),
	}

	// @todo Load partials
	// template.Must(site.Layout.ParseGlob(path.Join(partialsPath(), "*.html")))

	return result
}

// Build site
func (site *Site) Build() {
	// build nodes
	site.BuildActivities()
	site.BuildContact()
	site.BuildPages()
	site.BuildPosts()
	site.BuildHomepage()

	// dump errors
	site.ErrorCollector.Dump()
}

// Build activities node
func (site *Site) BuildActivities() {
	activitiesBuilder := NewActivitiesBuilder(site)
	activitiesBuilder.Load()
	activitiesBuilder.Generate(os.Stdout)
}

// Build contact node
func (site *Site) BuildContact() {
	contactBuilder := NewContactBuilder(site)
	contactBuilder.Load()
	contactBuilder.Generate(os.Stdout)
}

// Build pages nodes
func (site *Site) BuildPages() {
	pagesBuilder := NewPagesBuilder(site)
	pagesBuilder.Load()
	pagesBuilder.Generate(os.Stdout)
}

// Build posts nodes
func (site *Site) BuildPosts() {
	postsBuilder := NewPostsBuilder(site)
	postsBuilder.Load()
	postsBuilder.Generate(os.Stdout)
}

// Build homepage node
func (site *Site) BuildHomepage() {
	homepageBuilder := NewHomepageBuilder(site)
	homepageBuilder.Load()
	homepageBuilder.Generate(os.Stdout)
}

// Add an error when generating a node
func (site *Site) AddGenerationError(nodeKind string, err error) {
	step := fmt.Sprintf("Generating %s", nodeKind)
	site.ErrorCollector.AddError(step, err)
}
