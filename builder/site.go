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
	Model *models.Site

	// settings
	WorkingDir string
	OutputDir  string
	Theme      string
	UglyURL    bool

	// collectors
	ImageCollector *ImageCollector
	ErrorCollector *ErrorCollector

	// cache for #Layout method
	layout *template.Template

	// builders
	activitiesBuilder *ActivitiesBuilder
	contactBuilder    *ContactBuilder
	pagesBuilder      *PagesBuilder
	postsBuilder      *PostsBuilder
	homepageBuilder   *HomepageBuilder
}

func NewSite(siteId string) *Site {
	dbSession := models.NewDBSession()

	model := dbSession.FindSite(siteId)
	if model == nil {
		log.Fatalln("Can't find site with provided id")
	}

	result := &Site{
		Model: model,

		WorkingDir: viper.GetString("working_dir"),
		OutputDir:  viper.GetString("output_dir"),
		Theme:      viper.GetString("theme"),
		UglyURL:    viper.GetBool("ugly_url"),

		ImageCollector: NewImageCollector(),
		ErrorCollector: NewErrorCollector(),
	}

	return result
}

// Build site
func (site *Site) Build() {
	// init builders
	site.activitiesBuilder = NewActivitiesBuilder(site)
	site.activitiesBuilder.Load()

	site.contactBuilder = NewContactBuilder(site)
	site.contactBuilder.Load()

	site.pagesBuilder = NewPagesBuilder(site)
	site.pagesBuilder.Load()

	site.postsBuilder = NewPostsBuilder(site)
	site.postsBuilder.Load()

	site.homepageBuilder = NewHomepageBuilder(site)
	site.homepageBuilder.Load()

	// build nodes
	site.activitiesBuilder.Generate()
	site.contactBuilder.Generate()
	site.pagesBuilder.Generate()
	site.postsBuilder.Generate()
	site.homepageBuilder.Generate()

	// copy images
	site.CopyCollectedImages()

	// dump errors
	site.ErrorCollector.Dump()
}

// Copy images
func (site *Site) CopyCollectedImages() {
	errStep := "Copy images"

	imgDir := path.Join(site.GenDir(), IMAGES_DIR)

	// ensure img dir
	if err := site.EnsureDir(imgDir); err != nil {
		site.AddError(errStep, err)
		return
	}

	// copy images to img dir
	for _, imgKind := range site.ImageCollector.Images {
		derivative := models.DerivativeForKind(imgKind.Kind)
		srcFile := imgKind.Image.DerivativeFilePath(derivative)

		if err := site.CopyFile(srcFile, imgDir); err != nil {
			site.AddError(errStep, err)
		}
	}
}

// Copy file to given directory
func (site *Site) CopyFile(fromFilePath string, toDir string) error {
	// open source file
	src, err := os.Open(fromFilePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// open destination file
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

// Add an image to collector, and returns the URL for that image
func (site *Site) AddImage(img *models.Image, kind string) string {
	site.ImageCollector.AddImage(img, kind)

	// compute image URL
	// eg: /site_1/image_m.jpg => /img/image_m.jpg
	return "/" + path.Join(IMAGES_DIR, path.Base(img.DerivativeURL(models.DerivativeForKind(kind))))
}

// Add an error to collector
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
