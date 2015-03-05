package builder

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/fsync"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
)

const (
	IMAGES_DIR    = "img"
	THEMES_DIR    = "themes"
	TEMPLATES_DIR = "templates"
	PARTIALS_DIR  = "partials"
	ASSETS_DIR    = "assets"

	DEFAULT_THEME = "willy" // @todo FIXME !

	MAX_SLUG        = 50
	MAX_PAST_EVENTS = 5
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
	site   *models.Site
	config *SiteBuilderConfig

	images         []*models.Image
	errorCollector *ErrorCollector

	// all nodes slugs
	nodeSlugs map[string]bool

	// cache for #layout method
	masterLayout *template.Template

	// internal vars
	nodeBuilders map[string]NodeBuilder

	siteVars *SiteVars
	tplDir   string
}

type SiteBuilderConfig struct {
	WorkingDir string
	OutputDir  string
	Theme      string
	UglyURL    bool
	BaseURL    string
}

func NewSiteBuilder(site *models.Site, config *SiteBuilderConfig) *SiteBuilder {
	result := &SiteBuilder{
		site:   site,
		config: config,

		nodeSlugs: make(map[string]bool),

		errorCollector: NewErrorCollector(),

		nodeBuilders: make(map[string]NodeBuilder),
	}

	result.initBuilders()
	result.setTemplatesDir()

	return result
}

// Theme used by builder
func (builder *SiteBuilder) Theme() string {
	result := builder.config.Theme
	if result == "" {
		result = builder.site.Theme
	}

	if result == "" {
		result = DEFAULT_THEME
	}

	return result
}

// Are we building site with ugly urls ?
func (builder *SiteBuilder) UglyUrl() bool {
	if builder.config.UglyURL {
		return true
	} else {
		return builder.site.UglyURL
	}
}

// Build site
func (builder *SiteBuilder) Build() {
	// load nodes
	if builder.loadNodes(); builder.HaveError() {
		return
	}

	// compute site variables
	if builder.fillSiteVars(); builder.HaveError() {
		return
	}

	// sync nodes
	builder.syncNodes()

	// sync images
	builder.syncImages()

	// sync assets
	builder.syncAssets()
}

// Initialize builders
func (builder *SiteBuilder) initBuilders() {
	for name, initializer := range registeredNodeBuilders {
		builder.nodeBuilders[name] = initializer(builder)
	}
}

// Find templates dir
func (builder *SiteBuilder) setTemplatesDir() {
	dirPath := path.Join(builder.themeDir(), TEMPLATES_DIR)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// no /templates subdir found
		builder.tplDir = builder.themeDir()
	} else {
		builder.tplDir = dirPath
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

// Collect node slugs, and returns a new one if provided slug is already taken
func (builder *SiteBuilder) addNodeSlug(slug string) string {
	result := slug

	i := 1
	for builder.nodeSlugs[result] {
		result = fmt.Sprintf("%s-%d", slug, i)
		i++
	}

	builder.nodeSlugs[result] = true

	return result
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

// Returns all activities vars (to display in template)
func (builder *SiteBuilder) activitiesVars() []*ActivityVars {
	nodeBuilder := builder.nodeBuilder(KIND_ACTIVITIES)

	activities, ok := nodeBuilder.Data("activities").([]*ActivityVars)
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
func (builder *SiteBuilder) syncNodes() {
	errStep := "Sync nodes"

	allFiles := make(map[string]bool)
	allDirs := make(map[string]bool)

	// generate nodes
	for _, nodeBuilder := range builder.nodeBuilders {
		filePaths := nodeBuilder.Generate()
		for filePath, _ := range filePaths {
			log.Printf("Generated node: %s", filePath)
			allFiles[filePath] = true

			relativePath := path.Dir(strings.TrimPrefix(filePath, builder.GenDir()))

			destDir := builder.GenDir()
			for _, pathPart := range strings.Split(relativePath, "/") {
				if pathPart != "" {
					destDir = path.Join(destDir, pathPart)
					allDirs[destDir] = true
				}
			}
		}
	}

	// delete deprecated nodes
	filesToDelete := make(map[string]bool)

	// ignore assets and img dirs
	ignoreDirs := []string{
		path.Join(builder.GenDir(), "assets"),
		path.Join(builder.GenDir(), "img"),
	}

	err := filepath.Walk(builder.GenDir(), func(path string, f os.FileInfo, err error) error {
		if (path != builder.GenDir()) && !utils.HasOnePrefix(path, ignoreDirs) {
			if (f.IsDir() && !allDirs[path]) || (!f.IsDir() && !allFiles[path]) {
				filesToDelete[path] = true
			}
		}
		return nil
	})

	if err != nil {
		builder.addError(errStep, err)
		return
	}

	destFs := new(afero.OsFs)
	for filePath, _ := range filesToDelete {
		log.Printf("Deleting: %s", filePath)

		if err := destFs.RemoveAll(filePath); err != nil {
			builder.addError(errStep, err)
		}
	}
}

// Copy images
func (builder *SiteBuilder) syncImages() {
	errStep := "Sync images"

	imgDir := builder.genImagesDir()

	// ensure img dir
	if err := builder.ensureDir(imgDir); err != nil {
		builder.addError(errStep, err)
		return
	}

	// copy images to img dir
	sourceFiles := make(map[string]bool)

	var files []string
	for _, image := range builder.images {
		for _, derivative := range models.Derivatives {
			filePath := image.DerivativeFilePath(derivative)

			sourceFiles[path.Base(filePath)] = true

			files = append(files, filePath)
		}
	}

	if len(files) > 0 {
		fsync.SyncTo(imgDir, files...)
	}

	// delete deprecated images
	destFs := new(afero.OsFs)
	destfiles, err := afero.ReadDir(imgDir, destFs)
	if err != nil {
		builder.addError(errStep, err)
		return
	}

	for _, destfile := range destfiles {
		if !sourceFiles[destfile.Name()] {
			log.Printf("Deleting deprecated image: %s", destfile.Name())

			if err := destFs.RemoveAll(path.Join(imgDir, destfile.Name())); err != nil {
				builder.addError(errStep, err)
			}
		}
	}
}

// Copy theme assets
func (builder *SiteBuilder) syncAssets() error {
	syncer := fsync.NewSyncer()
	syncer.Delete = true

	return syncer.Sync(builder.genAssetsDir(), builder.themeAssetsDir())
}

// Collect image, and returns the template vars for that image
func (builder *SiteBuilder) addImage(img *models.Image) *ImageVars {
	builder.images = append(builder.images, img)

	return NewImageVars(img, builder.config.BaseURL)
}

// Check if builder have error
func (builder *SiteBuilder) HaveError() bool {
	return builder.errorCollector.ErrorsNb > 0
}

// Add an error to collector
func (builder *SiteBuilder) addError(step string, err error) {
	builder.errorCollector.addError(step, err)
}

// Add an error when generating a node
func (builder *SiteBuilder) addNodeBuilderError(nodeKind string, err error) {
	step := fmt.Sprintf("Building %s", nodeKind)
	builder.addError(step, err)
}

// Dump errors
func (builder *SiteBuilder) DumpErrors() {
	builder.errorCollector.dump()
}

// Computes theme directory
func (builder *SiteBuilder) themeDir() string {
	return path.Join(builder.config.WorkingDir, THEMES_DIR, builder.Theme())
}

// Computes theme templates directory
func (builder *SiteBuilder) themeTemplatesDir() string {
	if builder.tplDir == "" {
		panic("Templates directory not set")
	}

	return builder.tplDir
}

// Computes theme assets directory
func (builder *SiteBuilder) themeAssetsDir() string {
	return path.Join(builder.themeDir(), ASSETS_DIR)
}

// Computes directory where site is generated
func (builder *SiteBuilder) GenDir() string {
	return path.Join(builder.config.WorkingDir, builder.config.OutputDir)
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
	return path.Join(builder.themeTemplatesDir(), fmt.Sprintf("%s.html", tplName))
}

// Returns partials directory path
func (builder *SiteBuilder) partialsPath() string {
	return path.Join(builder.themeTemplatesDir(), PARTIALS_DIR)
}

func (builder *SiteBuilder) setupLayout() *template.Template {
	errStep := "template init"

	// parse layout
	result, err := template.ParseFiles(builder.templatePath("layout"))
	if err != nil {
		builder.addError(errStep, err)
		return nil
	}

	// setup FuncMap
	result.Funcs(builder.FuncMap())

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
					_, err := result.New(tplName).Parse(string(binData))
					if err != nil {
						builder.addError(errStep, err)
					}
				}
			}
		}
	}

	return result
}

// Get master layout template
func (builder *SiteBuilder) layout() *template.Template {
	if builder.masterLayout == nil {
		builder.masterLayout = builder.setupLayout()
	}

	return builder.masterLayout
}

// Dump templates
func (builder *SiteBuilder) DumpLayout() {
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
