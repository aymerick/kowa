package builder

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
)

type Site struct {
	Model          *models.Site
	ErrorCollector *ErrorCollector

	WorkingDir string
	OutputDir  string
	Theme      string

	layout *template.Template
}

func NewSite(siteId string) *Site {
	dbSession := models.NewDBSession()

	model := dbSession.FindSite(siteId)
	if model == nil {
		log.Fatalln("Can't find site with provided id")
	}

	result := &Site{
		Model:          model,
		ErrorCollector: NewErrorCollector(),

		WorkingDir: viper.GetString("working_dir"),
		OutputDir:  viper.GetString("output_dir"),
		Theme:      viper.GetString("theme"),
	}

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

// Computes directory where site is generated
func (site *Site) GenDir() string {
	return path.Join(site.WorkingDir, site.OutputDir)
}

// Prune directories for given absolute file path
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
	return path.Join(site.GenDir(), relativePath)
}

// Get master layout template
func (site *Site) Layout() *template.Template {
	if site.layout != nil {
		return site.layout
	} else {
		site.layout = template.Must(template.ParseFiles(site.TemplatePath("layout")))

		// @todo Load partials
		// template.Must(site.layout.ParseGlob(path.Join(site.PartialsPath(), "*.html")))

		return site.layout
	}
}

// Compute template path for given template name
func (site *Site) TemplatePath(tplName string) string {
	return path.Join(site.WorkingDir, "themes", site.Theme, fmt.Sprintf("%s.html", tplName))
}

// Returns partials directory path
func (site *Site) PartialsPath() string {
	return path.Join(site.WorkingDir, "themes", site.Theme, "partials")
}
