package builder

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/spf13/afero"
	"github.com/spf13/fsync"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/themes"
	"github.com/aymerick/raymond"
)

const (
	// static site
	assetsDir       = "assets"
	imagesDir       = "img"
	filesDir        = "files"
	faviconFilename = "favicon.png"

	// themes
	templatesDir = "templates"
	partialsDir  = "partials"

	maxSlug       = 50
	maxPastEvents = 5
)

var generatedPaths = []string{assetsDir, imagesDir, filesDir, faviconFilename}

var registeredNodeBuilders = make(map[string]func(*SiteBuilder) NodeBuilder)

// RegisterNodeBuilder registers a NodeBuilder
func RegisterNodeBuilder(name string, initializer func(*SiteBuilder) NodeBuilder) {
	if _, exists := registeredNodeBuilders[name]; exists {
		panic(fmt.Sprintf("Builder initializer already registered: %s", name))
	}

	registeredNodeBuilders[name] = initializer
}

// SiteBuilder builds an entire site
type SiteBuilder struct {
	site  *models.Site
	theme *themes.Theme

	images         []*models.Image
	files          []*models.File
	errorCollector *ErrorCollector

	// all nodes slugs
	nodeSlugs map[string]bool

	// cache for #layout method
	masterLayout *raymond.Template

	// internal vars
	nodeBuilders map[string]NodeBuilder

	siteVars *SiteVars
}

// NewSiteBuilder instanciates a new SiteBuilder
func NewSiteBuilder(site *models.Site) *SiteBuilder {
	result := &SiteBuilder{
		site:  site,
		theme: themes.New(site.Theme),

		nodeSlugs:      make(map[string]bool),
		errorCollector: NewErrorCollector(),
		nodeBuilders:   make(map[string]NodeBuilder),
	}

	result.initBuilders()

	return result
}

// Build executes site building
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
	builder.syncFiles(builder.genImagesDir(), builder.imagesToSync)

	// sync files
	builder.syncFiles(builder.genFilesDir(), builder.filesToSync)

	// sync assets
	builder.syncAssets()

	// compile SASS files into CSS
	builder.buildSass()

	// sync favicon
	builder.syncFavicon()
}

// OutputDir returns path to output directory
func (builder *SiteBuilder) OutputDir() string {
	return path.Join(viper.GetString("output_dir"), builder.site.BuildDir())
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
	nodeBuilder := builder.nodeBuilder(kindActivities)

	activities, ok := nodeBuilder.Data("activities").([]*ActivityVars)
	if !ok {
		panic("This should never happen")
	}

	return activities
}

// Return site base path
func (builder *SiteBuilder) basePath() string {
	u, err := url.Parse(builder.site.BaseUrl())
	if err != nil {
		return ""
	}

	path := strings.TrimSuffix(u.Path, "/")

	return path
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

		for filePath := range filePaths {
			log.Printf("Generated node: %+q", filePath)
			allFiles[filePath] = true

			relativePath := path.Dir(strings.TrimPrefix(filePath, builder.OutputDir()))

			destDir := builder.OutputDir()
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

	// ignore generated dirs
	var ignoreDirs []string

	for _, genPath := range generatedPaths {
		ignoreDirs = append(ignoreDirs, path.Join(builder.OutputDir(), genPath))
	}

	err := filepath.Walk(builder.OutputDir(), func(path string, f os.FileInfo, err error) error {
		if (path != builder.OutputDir()) && !helpers.HasOnePrefix(path, ignoreDirs) {
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
	for filePath := range filesToDelete {
		log.Printf("Deleting: %+q", filePath)

		if err := destFs.RemoveAll(filePath); err != nil {
			builder.addError(errStep, err)
		}
	}
}

func (builder *SiteBuilder) imagesToSync() ([]string, map[string]bool) {
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

	return files, sourceFiles
}

func (builder *SiteBuilder) filesToSync() ([]string, map[string]bool) {
	sourceFiles := make(map[string]bool)

	var files []string
	for _, file := range builder.files {
		filePath := file.FilePath()

		if !sourceFiles[path.Base(filePath)] {
			files = append(files, filePath)

			sourceFiles[path.Base(filePath)] = true
		}
	}

	return files, sourceFiles
}

// Sync files
func (builder *SiteBuilder) syncFiles(destDir string, grabber func() ([]string, map[string]bool)) {
	errStep := "Sync " + destDir

	// ensure dir
	if err := builder.ensureDir(destDir); err != nil {
		builder.addError(errStep, err)
		return
	}

	// grab files to sync
	files, sourceFiles := grabber()

	if len(files) > 0 {
		log.Printf("Syncing %d files", len(files))

		if err := fsync.SyncTo(destDir, files...); err != nil {
			builder.addError(errStep, err)
		}
	}

	// delete deprecated files
	destFs := new(afero.OsFs)
	destfiles, err := afero.ReadDir(destDir, destFs)
	if err != nil {
		builder.addError(errStep, err)
		return
	}

	for _, destfile := range destfiles {
		if !sourceFiles[destfile.Name()] {
			log.Printf("Deleting deprecated file: %s", destfile.Name())

			if err := destFs.RemoveAll(path.Join(destDir, destfile.Name())); err != nil {
				builder.addError(errStep, err)
			}
		}
	}
}

// Compile theme SASS files into CSS
func (builder *SiteBuilder) buildSass() {
	if !builder.theme.HaveSass() {
		// no sass dir
		return
	}

	errStep := "Build SASS"

	log.Printf("Compiling SASS file(s)")

	// computes sass vars
	sassVars, err := builder.computeSassVars()
	if err != nil {
		builder.addError(errStep, err)
		return
	}

	// compile sass files
	if err := builder.theme.SassBuild(sassVars, builder.genAssetsDir()); err != nil {
		builder.addError(errStep, err)
	}
}

// computes SASS variables for current site and theme
func (builder *SiteBuilder) computeSassVars() (string, error) {
	result := ""

	// get theme variables
	themeVars, err := builder.theme.SassVars()
	if err != nil {
		return "", err
	}

	// fetch sass variables from database
	siteVars := builder.siteSassVariables()

	// computes all variables
	for name, value := range themeVars {
		result += "$" + name + ": "

		if siteVars[name] != "" {
			// custom value
			value = siteVars[name]
		}

		result += value + ";\n"
	}

	return result, nil
}

// Fetches SASS variables from database
func (builder *SiteBuilder) siteSassVariables() map[string]string {
	result := map[string]string{}

	// @todo FIXME !
	if settings := builder.site.ThemeSettings[builder.site.Theme]; settings != nil {
		for _, sassVar := range settings.Sass {
			result[sassVar.Name] = sassVar.Value
		}
	}

	return result
}

// Copy theme assets
func (builder *SiteBuilder) syncAssets() {
	errStep := "Sync assets"

	syncer := fsync.NewSyncer()
	syncer.Delete = true
	syncer.NoTimes = true

	if err := syncer.Sync(builder.genAssetsDir(), builder.theme.AssetsDir); err != nil {
		builder.addError(errStep, err)
	}
}

// Copy or generate favicon
func (builder *SiteBuilder) syncFavicon() {
	errStep := "Sync favicon"

	faviconPath := path.Join(builder.OutputDir(), faviconFilename)

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

	return NewImageVars(img, builder.basePath(), builder.site.BaseUrl())
}

// Collect file, and returns the URL for that file
func (builder *SiteBuilder) addFile(file *models.File) string {
	builder.files = append(builder.files, file)

	return builder.site.BaseUrl() + path.Join("/", filesDir, file.Path)
}

// HaveError returns true if builder have error
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

// DumpErrors displays collected errors
func (builder *SiteBuilder) DumpErrors() {
	builder.errorCollector.dump()
}

// Computes directory where images are copied
func (builder *SiteBuilder) genImagesDir() string {
	return path.Join(builder.OutputDir(), imagesDir)
}

// Computes directory where files are copied
func (builder *SiteBuilder) genFilesDir() string {
	return path.Join(builder.OutputDir(), filesDir)
}

// Computes directory where assets are copied
func (builder *SiteBuilder) genAssetsDir() string {
	return path.Join(builder.OutputDir(), assetsDir)
}

// @todo In production: load all layout files only once on startup, then for each builder instance:
//       clone layout and register helpers
func (builder *SiteBuilder) setupLayout() *raymond.Template {
	errStep := "Layout setup"

	// parse layout
	result, err := raymond.ParseFile(builder.theme.Template("layout"))
	if err != nil {
		builder.addError(errStep, err)
		return nil
	}

	// register helpers
	result.RegisterHelpers(builder.helpers())

	// register partials
	filePaths, err := builder.theme.Partials()
	if err != nil {
		builder.addError(errStep, err)
		return nil
	}

	result.RegisterPartialFiles(filePaths...)

	return result
}

// Get master layout template
func (builder *SiteBuilder) layout() *raymond.Template {
	if builder.masterLayout == nil {
		builder.masterLayout = builder.setupLayout()
	}

	return builder.masterLayout
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
	return path.Join(builder.OutputDir(), relativePath)
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
