package builder

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/fsync"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
)

const (
	IMAGES_DIR   = "img"
	THEMES_DIR   = "themes"
	PARTIALS_DIR = "partials"
	ASSETS_DIR   = "assets"
)

var registeredNodeBuilders = make(map[string]func(*SiteBuilder) NodeBuilder)

// Registers a node builder
func RegisterNodeBuilder(name string, initializer func(*SiteBuilder) NodeBuilder) {
	if _, exists := registeredNodeBuilders[name]; exists {
		panic(fmt.Sprintf("Builder initializer already registered: %s", name))
	}

	registeredNodeBuilders[name] = initializer
}

type SiteBuilder struct {
	// settings
	workingDir string
	outputDir  string
	theme      string
	uglyURL    bool

	// collectors
	imageCollector *ImageCollector
	errorCollector *ErrorCollector

	// cache for #layout method
	masterLayout *template.Template

	site         *models.Site
	siteVars     *SiteVars
	nodeBuilders map[string]NodeBuilder
}

func NewSiteBuilder(siteId string) *SiteBuilder {
	dbSession := models.NewDBSession()

	site := dbSession.FindSite(siteId)
	if site == nil {
		log.Fatalln("Can't find site with provided id")
	}

	result := &SiteBuilder{
		site: site,

		workingDir: viper.GetString("working_dir"),
		outputDir:  viper.GetString("output_dir"),
		theme:      viper.GetString("theme"),
		uglyURL:    viper.GetBool("ugly_url"),

		imageCollector: NewImageCollector(),
		errorCollector: NewErrorCollector(),

		nodeBuilders: make(map[string]NodeBuilder),
	}

	result.initBuilders()

	return result
}

// Build site
func (builder *SiteBuilder) Build() {
	// load nodes
	builder.loadNodes()

	// compute site variables
	builder.fillSiteVars()

	// generate nodes
	builder.generateNodes()

	// copy images
	builder.copyCollectedImages()

	// copy assets
	builder.copyAssets()

	// check errors
	if builder.haveError() {
		builder.dumpErrors()
		builder.dumpLayout()
	}
}

// Initialize builders
func (builder *SiteBuilder) initBuilders() {
	for name, initializer := range registeredNodeBuilders {
		builder.nodeBuilders[name] = initializer(builder)
	}
}

// Get given node builder
func (builder *SiteBuilder) nodeBuilder(name string) NodeBuilder {
	return builder.nodeBuilders[name]
}

// Load nodes
func (builder *SiteBuilder) loadNodes() {
	for _, nodeBuilder := range builder.nodeBuilders {
		nodeBuilder.Load()
	}
}

// Returns nodes to display in navigation bar
func (builder *SiteBuilder) navBarNodes() []*Node {
	result := []*Node{}

	for _, nodeBuilder := range builder.nodeBuilders {
		nodes := nodeBuilder.NavBarNodes()
		if len(nodes) > 0 {
			result = append(result, nodes...)
		}
	}

	return result
}

// Returns all activities contents (to display in template)
func (builder *SiteBuilder) activitiesContents() []*ActivityContent {
	nodeBuilder := builder.nodeBuilder(KIND_ACTIVITIES)

	activities, ok := nodeBuilder.Data("activities").([]*ActivityContent)
	if !ok {
		panic("This should never happen")
	}

	return activities
}

// Fill site variables
func (builder *SiteBuilder) fillSiteVars() {
	builder.siteVars = NewSiteVars(builder)
	builder.siteVars.fill()
}

// Generate nodes
func (builder *SiteBuilder) generateNodes() {
	for _, nodeBuilder := range builder.nodeBuilders {
		nodeBuilder.Generate()
	}
}

// Copy images
func (builder *SiteBuilder) copyCollectedImages() {
	errStep := "Copy images"

	imgDir := builder.genImagesDir()

	// ensure img dir
	if err := builder.ensureDir(imgDir); err != nil {
		builder.addError(errStep, err)
		return
	}

	// copy images to img dir
	for _, imgKind := range builder.imageCollector.Images {
		derivative := models.DerivativeForKind(imgKind.Kind)
		srcFile := imgKind.Image.DerivativeFilePath(derivative)

		if err := builder.copyFile(srcFile, imgDir); err != nil {
			builder.addError(errStep, err)
		}
	}
}

// Copy theme assets
func (builder *SiteBuilder) copyAssets() error {
	syncer := fsync.NewSyncer()
	syncer.SrcFs = new(afero.OsFs)
	syncer.DestFs = new(afero.OsFs)

	return syncer.Sync(builder.genAssetsDir(), builder.themeAssetsDir())
}

// Add an image to collector, and returns the URL for that image
func (builder *SiteBuilder) addImage(img *models.Image, kind string) string {
	builder.imageCollector.addImage(img, kind)

	// compute image URL
	// eg: /site_1/image_m.jpg => /img/image_m.jpg
	return "/" + path.Join(IMAGES_DIR, path.Base(img.DerivativeURL(models.DerivativeForKind(kind))))
}

// Check if builder have error
func (builder *SiteBuilder) haveError() bool {
	return builder.errorCollector.ErrorsNb > 0
}

// Add an error to collector
func (builder *SiteBuilder) addError(step string, err error) {
	builder.errorCollector.addError(step, err)
}

// Add an error when generating a node
func (builder *SiteBuilder) addGenError(nodeKind string, err error) {
	step := fmt.Sprintf("Generating %s", nodeKind)
	builder.addError(step, err)
}

// Dump errors
func (builder *SiteBuilder) dumpErrors() {
	builder.errorCollector.dump()
}

// Computes theme directory
func (builder *SiteBuilder) themeDir() string {
	return path.Join(builder.workingDir, THEMES_DIR, builder.theme)
}

// Computes theme assets directory
func (builder *SiteBuilder) themeAssetsDir() string {
	return path.Join(builder.themeDir(), ASSETS_DIR)
}

// Computes directory where site is generated
func (builder *SiteBuilder) GenDir() string {
	return path.Join(builder.workingDir, builder.outputDir)
}

// Computes directory where images are copied
func (builder *SiteBuilder) genImagesDir() string {
	return path.Join(builder.GenDir(), IMAGES_DIR)
}

// Computes directory where assets are copied
func (builder *SiteBuilder) genAssetsDir() string {
	return path.Join(builder.GenDir(), ASSETS_DIR)
}

// Compute template path for given template name
func (builder *SiteBuilder) templatePath(tplName string) string {
	return path.Join(builder.themeDir(), fmt.Sprintf("%s.html", tplName))
}

// Returns partials directory path
func (builder *SiteBuilder) partialsPath() string {
	return path.Join(builder.themeDir(), PARTIALS_DIR)
}

// Get master layout template
func (builder *SiteBuilder) layout() *template.Template {
	if builder.masterLayout != nil {
		return builder.masterLayout
	} else {
		errStep := "template init"

		// parse layout
		layout, err := template.ParseFiles(builder.templatePath("layout"))
		if err != nil {
			builder.addError(errStep, err)
			return nil
		}

		builder.masterLayout = layout

		// load partials
		partialDir := builder.partialsPath()

		files, err := ioutil.ReadDir(partialDir)
		if err != nil && err != os.ErrExist {
			builder.addError(errStep, err)
		} else {
			for _, file := range files {
				fileName := file.Name()

				if !file.IsDir() && strings.HasSuffix(fileName, ".html") {
					filePath := path.Join(partialDir, fileName)

					// read partial
					binData, err := ioutil.ReadFile(filePath)
					if err != nil {
						builder.addError(errStep, err)
					} else {
						// eg: partials/navbar
						tplName := fmt.Sprintf("%s/%s", PARTIALS_DIR, utils.FileBase(fileName))

						// add partial to layout
						_, err := builder.masterLayout.New(tplName).Parse(string(binData))
						if err != nil {
							builder.addError(errStep, err)
						}
					}
				}
			}
		}

		return builder.masterLayout
	}
}

// Dump templates
func (builder *SiteBuilder) dumpLayout() {
	log.Printf("Layout templates:")
	for _, tpl := range builder.masterLayout.Templates() {
		log.Printf("  -> %s", tpl.Name())
	}
}

// Prune directories for given absolute dir path
func (builder *SiteBuilder) ensureDir(dirPath string) error {
	// log.Printf("[DBG] Creating dir: %s", dirPath)

	err := os.MkdirAll(dirPath, 0777)
	if err != nil && err != os.ErrExist {
		return err
	}

	return err
}

// Prune directories for given absolute file path
func (builder *SiteBuilder) ensureFileDir(osPath string) error {
	return builder.ensureDir(path.Dir(osPath))
}

// Computes local file path for given URL
func (builder *SiteBuilder) filePath(fullUrl string) string {
	return path.Join(builder.GenDir(), fullUrl)
}

// Copy file to given directory
func (builder *SiteBuilder) copyFile(fromFilePath string, toDir string) error {
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
