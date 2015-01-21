package builder

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Site struct {
	Layout         *template.Template
	ErrorCollector *ErrorCollector

	GenDir string
}

func NewSite() *Site {
	result := &Site{
		Layout:         template.Must(template.ParseFiles(templatePath("layout"))),
		ErrorCollector: NewErrorCollector(),
		GenDir:         path.Join(viper.GetString("working_dir"), viper.GetString("output_dir")),
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
	activitiesBuilder.Generate()
}

// Build contact node
func (site *Site) BuildContact() {
	contactBuilder := NewContactBuilder(site)
	contactBuilder.Load()
	contactBuilder.Generate()
}

// Build pages nodes
func (site *Site) BuildPages() {
	pagesBuilder := NewPagesBuilder(site)
	pagesBuilder.Load()
	pagesBuilder.Generate()
}

// Build posts nodes
func (site *Site) BuildPosts() {
	postsBuilder := NewPostsBuilder(site)
	postsBuilder.Load()
	postsBuilder.Generate()
}

// Build homepage node
func (site *Site) BuildHomepage() {
	homepageBuilder := NewHomepageBuilder(site)
	homepageBuilder.Load()
	homepageBuilder.Generate()
}

// Add an error when generating a node
func (site *Site) AddGenError(nodeKind string, err error) {
	step := fmt.Sprintf("Generating %s", nodeKind)
	site.ErrorCollector.AddError(step, err)
}

func (site *Site) EnsureFileDir(osPath string) error {
	fileDir := path.Dir(osPath)

	log.Printf("[DBG] Creating dir: %s", fileDir)

	err := os.MkdirAll(fileDir, 0777)
	if err != nil && err != os.ErrExist {
		return err
	}

	return err
}

// Computes an absolute file path
func (site *Site) FilePath(relativePath string) string {
	return path.Join(site.GenDir, relativePath)
}
