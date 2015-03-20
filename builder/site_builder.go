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

	"github.com/disintegration/imaging"
	"github.com/spf13/afero"
	"github.com/spf13/fsync"

	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/models"
)

const (
	IMAGES_DIR    = "img"
	TEMPLATES_DIR = "templates"
	PARTIALS_DIR  = "partials"
	ASSETS_DIR    = "assets"

	FAVICON_FILENAME = "favicon.png"

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
	ThemesDir string
	OutputDir string
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

func (builder *SiteBuilder) Config() *SiteBuilderConfig {
	return builder.config
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

	// sync favicon
	builder.syncFavicon()
}

// Initialize builders
func (builder *SiteBuilder) initBuilders() {
	for name, initializer := range registeredNodeBuilders {
		builder.nodeBuilders[name] = initializer(builder)
	}
}

// Set templates dir
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

// Return site base path
func (builder *SiteBuilder) basePath() string {
	_, result := path.Split(builder.site.BaseUrl)
	return path.Join("/", result)
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
	// @todo Use go routines and channels
	for _, nodeBuilder := range builder.nodeBuilders {
		filePaths := nodeBuilder.Generate()

		for filePath, _ := range filePaths {
			log.Printf("Generated node: %+q", filePath)
			allFiles[filePath] = true

			relativePath := path.Dir(strings.TrimPrefix(filePath, builder.config.OutputDir))

			destDir := builder.config.OutputDir
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
		path.Join(builder.config.OutputDir, "assets"),
		path.Join(builder.config.OutputDir, "img"),
	}

	err := filepath.Walk(builder.config.OutputDir, func(path string, f os.FileInfo, err error) error {
		if (path != builder.config.OutputDir) && !helpers.HasOnePrefix(path, ignoreDirs) {
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
		log.Printf("Deleting: %+q", filePath)

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

			if !sourceFiles[path.Base(filePath)] {
				files = append(files, filePath)

				sourceFiles[path.Base(filePath)] = true
			}
		}
	}

	if len(files) > 0 {
		log.Printf("Syncing %d images", len(files))

		if err := fsync.SyncTo(imgDir, files...); err != nil {
			builder.addError(errStep, err)
		}
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
func (builder *SiteBuilder) syncAssets() {
	errStep := "Sync assets"

	syncer := fsync.NewSyncer()
	syncer.Delete = true
	syncer.NoTimes = true

	if err := syncer.Sync(builder.genAssetsDir(), builder.themeAssetsDir()); err != nil {
		builder.addError(errStep, err)
	}
}

// Copy or generate favicon
func (builder *SiteBuilder) syncFavicon() {
	errStep := "Sync favicon"

	faviconPath := path.Join(builder.config.OutputDir, FAVICON_FILENAME)

	if img := builder.site.FindFavicon(); img != nil {
		log.Printf("Generating favicon")

		// build 16x16 favicon
		favicon := imaging.Thumbnail(*img.Original(), 16, 16, imaging.Lanczos)

		// save favicon
		if err := imaging.Save(favicon, faviconPath); err != nil {
			builder.addError(errStep, err)
		}
	} else {
		// delete deprecated favicon
		if _, err := os.Stat(faviconPath); !os.IsNotExist(err) {
			log.Printf("Stat: %v", err)

			if errRem := os.Remove(faviconPath); errRem != nil {
				builder.addError(errStep, errRem)
			}
		}
	}
}

// Collect image, and returns the template vars for that image
func (builder *SiteBuilder) addImage(img *models.Image) *ImageVars {
	builder.images = append(builder.images, img)

	return NewImageVars(img, builder.basePath(), builder.site.BaseUrl)
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
	return path.Join(builder.config.ThemesDir, builder.site.Theme)
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

// Computes directory where images are copied
func (builder *SiteBuilder) genImagesDir() string {
	return path.Join(builder.config.OutputDir, IMAGES_DIR)
}

// Computes directory where assets are copied
func (builder *SiteBuilder) genAssetsDir() string {
	return path.Join(builder.config.OutputDir, ASSETS_DIR)
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
					tplName := fmt.Sprintf("%s/%s", PARTIALS_DIR, helpers.FileBase(fileName))

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

// Computes local file path for given relative path
func (builder *SiteBuilder) filePath(relativePath string) string {
	return path.Join(builder.config.OutputDir, relativePath)
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
