package builder

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"

	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
)

const (
	IMAGES_DIR = "img"
)

type Site struct {
	Model          *models.Site
	ImageCollector *ImageCollector
	ErrorCollector *ErrorCollector

	WorkingDir string
	OutputDir  string
	Theme      string
	UglyURL    bool

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
		ImageCollector: NewImageCollector(),
		ErrorCollector: NewErrorCollector(),

		WorkingDir: viper.GetString("working_dir"),
		OutputDir:  viper.GetString("output_dir"),
		Theme:      viper.GetString("theme"),
		UglyURL:    viper.GetBool("ugly_url"),
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

	// copy images
	site.CopyCollectedImages()

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

// Copy images
func (site *Site) CopyCollectedImages() {
	errStep := "Copy images"
	imgDir := path.Join(site.GenDir(), IMAGES_DIR)

	if err := site.EnsureDir(imgDir); err != nil {
		site.AddError(errStep, err)
		return
	}

	for _, imgKind := range site.ImageCollector.Images {
		// Copy medium image
		derivative := models.DerivativeForKind(imgKind.Kind)
		srcFile := imgKind.Image.DerivativeFilePath(derivative)

		if err := site.CopyFile(srcFile, imgDir); err != nil {
			site.AddError(errStep, err)
		}
	}
}

// Copy a file to given directory
func (site *Site) CopyFile(fromFilePath string, toDir string) error {
	// open source
	src, err := os.Open(fromFilePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// open destination
	dstFilePath := path.Join(toDir, path.Base(fromFilePath))

	dst, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// Add an image to process, and returns the URL for that image
func (site *Site) AddImage(img *models.Image, kind string) string {
	site.ImageCollector.AddImage(img, kind)

	// fix image URL
	// eg: /site_1/image_m.jpg => /img/image_m.jpg
	return "/" + path.Join(IMAGES_DIR, path.Base(img.DerivativeURL(models.DerivativeForKind(kind))))
}

// Add an error
func (site *Site) AddError(step string, err error) {
	site.ErrorCollector.AddError(step, err)
}

// Add an error when generating a node
func (site *Site) AddGenError(nodeKind string, err error) {
	step := fmt.Sprintf("Generating %s", nodeKind)
	site.AddError(step, err)
}

// Computes directory where site is generated
func (site *Site) GenDir() string {
	return path.Join(site.WorkingDir, site.OutputDir)
}

// Prune directories for given absolute dir path
func (site *Site) EnsureDir(dirPath string) error {
	// log.Printf("[DBG] Creating dir: %s", dirPath)

	err := os.MkdirAll(dirPath, 0777)
	if err != nil && err != os.ErrExist {
		return err
	}

	return err
}

// Prune directories for given absolute file path
func (site *Site) EnsureFileDir(osPath string) error {
	return site.EnsureDir(path.Dir(osPath))
}

// Computes local file path for given URL
func (site *Site) FilePath(fullUrl string) string {
	return path.Join(site.GenDir(), fullUrl)
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
