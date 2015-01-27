package builder

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"

	"github.com/spf13/afero"
	"github.com/spf13/fsync"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/models"
)

const (
	IMAGES_DIR   = "img"
	THEMES_DIR   = "themes"
	PARTIALS_DIR = "partials"
	ASSETS_DIR   = "assets"
)

var builderInitializers = make(map[string]func(*Site) NodeBuilder)

// Registers a builder
func RegisterBuilderInitializer(name string, initializer func(*Site) NodeBuilder) {
	if _, exists := builderInitializers[name]; exists {
		panic(fmt.Sprintf("Builder initializer already registered: %s", name))
	}

	builderInitializers[name] = initializer
}

type Site struct {
	// settings
	workingDir string
	outputDir  string
	theme      string
	uglyURL    bool

	// collectors
	imageCollector *ImageCollector
	errorCollector *ErrorCollector

	// cache for #Layout method
	layout *template.Template

	model    *models.Site
	builders map[string]NodeBuilder
}

func NewSite(siteId string) *Site {
	dbSession := models.NewDBSession()

	model := dbSession.FindSite(siteId)
	if model == nil {
		log.Fatalln("Can't find site with provided id")
	}

	result := &Site{
		model: model,

		workingDir: viper.GetString("working_dir"),
		outputDir:  viper.GetString("output_dir"),
		theme:      viper.GetString("theme"),
		uglyURL:    viper.GetBool("ugly_url"),

		imageCollector: NewImageCollector(),
		errorCollector: NewErrorCollector(),

		builders: make(map[string]NodeBuilder),
	}

	result.initBuilders()

	return result
}

// Initialize builders
func (site *Site) initBuilders() {
	for name, initializer := range builderInitializers {
		site.builders[name] = initializer(site)
	}
}

// Build site
func (site *Site) Build() {
	// load nodes
	site.loadNodes()

	// generate nodes
	site.generateNodes()

	// copy images
	site.CopyCollectedImages()

	// copy assets
	site.CopyAssets()

	// dump errors
	site.errorCollector.Dump()
}

// Load nodes
func (site *Site) loadNodes() {
	for _, builder := range site.builders {
		builder.Load()
	}
}

// Generate nodes
func (site *Site) generateNodes() {
	for _, builder := range site.builders {
		builder.Generate()
	}
}

// Copy images
func (site *Site) CopyCollectedImages() {
	errStep := "Copy images"

	imgDir := site.GenImagesDir()

	// ensure img dir
	if err := site.EnsureDir(imgDir); err != nil {
		site.AddError(errStep, err)
		return
	}

	// copy images to img dir
	for _, imgKind := range site.imageCollector.Images {
		derivative := models.DerivativeForKind(imgKind.Kind)
		srcFile := imgKind.Image.DerivativeFilePath(derivative)

		if err := site.CopyFile(srcFile, imgDir); err != nil {
			site.AddError(errStep, err)
		}
	}
}

// Copy theme assets
func (site *Site) CopyAssets() error {
	syncer := fsync.NewSyncer()
	syncer.SrcFs = new(afero.OsFs)
	syncer.DestFs = new(afero.OsFs)

	return syncer.Sync(site.GenAssetsDir(), site.ThemeAssetsDir())
}

// Add an image to collector, and returns the URL for that image
func (site *Site) AddImage(img *models.Image, kind string) string {
	site.imageCollector.AddImage(img, kind)

	// compute image URL
	// eg: /site_1/image_m.jpg => /img/image_m.jpg
	return "/" + path.Join(IMAGES_DIR, path.Base(img.DerivativeURL(models.DerivativeForKind(kind))))
}

// Add an error to collector
func (site *Site) AddError(step string, err error) {
	site.errorCollector.AddError(step, err)
}

// Add an error when generating a node
func (site *Site) AddGenError(nodeKind string, err error) {
	step := fmt.Sprintf("Generating %s", nodeKind)
	site.AddError(step, err)
}

// Computes theme directory
func (site *Site) ThemeDir() string {
	return path.Join(site.workingDir, THEMES_DIR, site.theme)
}

// Computes theme assets directory
func (site *Site) ThemeAssetsDir() string {
	return path.Join(site.ThemeDir(), ASSETS_DIR)
}

// Computes directory where site is generated
func (site *Site) GenDir() string {
	return path.Join(site.workingDir, site.outputDir)
}

// Copmputes directory where images are copied
func (site *Site) GenImagesDir() string {
	return path.Join(site.GenDir(), IMAGES_DIR)
}

// Copmputes directory where assets are copied
func (site *Site) GenAssetsDir() string {
	return path.Join(site.GenDir(), ASSETS_DIR)
}

// Compute template path for given template name
func (site *Site) TemplatePath(tplName string) string {
	return path.Join(site.ThemeDir(), fmt.Sprintf("%s.html", tplName))
}

// Get master layout template
func (site *Site) Layout() *template.Template {
	if site.layout != nil {
		return site.layout
	} else {
		site.layout = template.Must(template.ParseFiles(site.TemplatePath("layout")))

		// Load partials
		template.Must(site.layout.ParseGlob(path.Join(site.PartialsPath(), "*.html")))

		// for _, tpl := range site.layout.Templates() {
		// 	log.Printf("Template: %s", tpl.Name())
		// }

		return site.layout
	}
}

// Returns partials directory path
func (site *Site) PartialsPath() string {
	return path.Join(site.ThemeDir(), PARTIALS_DIR)
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
