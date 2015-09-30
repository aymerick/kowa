package themes

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/aymerick/kowa/helpers"
)

const (
	// theme templates subdirectory
	templatesDir = "templates"

	// theme templates partials subdirectory
	partialsDir = "partials"

	// theme assets subdirectory
	assetsDir = "assets"

	// theme sass subdirectory
	sassDir = "sass"

	// default palette name
	defaultPalette = "Default"

	// configuration file name
	confFile = "theme.toml"
)

// Theme represents a theme
type Theme struct {
	ID string

	Dir          string
	AssetsDir    string
	TemplatesDir string
	PartialsDir  string
	SassDir      string
	Conf         *Conf

	partials  []string
	sassFiles []os.FileInfo
	sassVars  map[string]string
	palettes  []*Palette
}

// New instanciates a new Theme
func New(id string) *Theme {
	t := &Theme{
		ID:  id,
		Dir: path.Join(Dir(), id),

		sassVars: map[string]string{},
	}

	t.AssetsDir = path.Join(t.Dir, assetsDir)

	t.setTemplateDir()
	t.setPartialsDir()
	t.setSassDir()
	t.setConf()

	return t
}

// AssetExist returns true if asset at given relative path exists in theme
func (t *Theme) AssetExist(relativePath string) bool {
	_, err := os.Stat(path.Join(t.AssetsDir, relativePath))

	return !os.IsNotExist(err)
}

// Template returns template path
func (t *Theme) Template(id string) string {
	return path.Join(t.TemplatesDir, fmt.Sprintf("%s.hbs", id))
}

// Partials returns an array of partials paths
func (t *Theme) Partials() ([]string, error) {
	// @todo Recompute on file change
	if len(t.partials) == 0 {
		files, err := ioutil.ReadDir(t.PartialsDir)
		if err != nil && err != os.ErrExist {
			return t.partials, err
		}

		for _, file := range files {
			fileName := file.Name()

			if !file.IsDir() && strings.HasSuffix(fileName, ".hbs") {
				t.partials = append(t.partials, path.Join(t.PartialsDir, fileName))
			}
		}
	}

	return t.partials, nil
}

// HaveSass returns true if theme have sass files
func (t *Theme) HaveSass() bool {
	return t.SassDir != ""
}

// SassFiles returns theme sass files infos
func (t *Theme) SassFiles() ([]os.FileInfo, error) {
	if t.SassDir == "" {
		// no sass dir
		return t.sassFiles, nil
	}

	// @todo Recompute on file change
	if len(t.sassFiles) == 0 {
		var err error
		if t.sassFiles, err = ioutil.ReadDir(t.SassDir); err != nil {
			return t.sassFiles, err
		}
	}

	return t.sassFiles, nil
}

// SassFile returns the absolute path to given sass file
func (t *Theme) SassFile(subPath string) string {
	return path.Join(t.SassDir, subPath)
}

// SassVars returns the sass variables found in theme
func (t *Theme) SassVars() (map[string]string, error) {
	// @todo Recompute on file change
	if len(t.sassVars) == 0 {
		// parse variables file from theme
		file, err := os.Open(path.Join(t.SassDir, "_variables.scss"))
		if err != nil {
			if os.IsNotExist(err) {
				// that's ok, theme have no variables file
				return t.sassVars, nil
			}
			return t.sassVars, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "$") {
				if pair := strings.SplitN(line, ":", 2); len(pair) == 2 {
					name := strings.TrimSpace(pair[0][1:len(pair[0])])
					value := strings.TrimSpace(pair[1])

					if value[len(value)-1] == ';' {
						value = strings.TrimSpace(value[0 : len(value)-1])
						t.sassVars[name] = value
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return t.sassVars, err
		}
	}

	return t.sassVars, nil
}

// SassBuild builds all sass files with given sass vars overwrites, and outputs result into given directory
func (t *Theme) SassBuild(sassVars string, output string) error {
	files, err := t.SassFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		sassPath := t.SassFile(file.Name())
		baseName := helpers.FileBase(sassPath)

		cssRelativePath := path.Join("css", helpers.FileBase(sassPath)+".css")

		// if CSS file exists in theme, do NOT overwrite it
		if t.AssetExist(cssRelativePath) {
			log.Printf("Skipping SASS file '%s' because that CSS file is present in theme: %s", sassPath, cssRelativePath)
		} else {
			outPath := path.Join(output, cssRelativePath)

			// skip directories and partials
			if strings.HasSuffix(sassPath, ".scss") && !file.IsDir() && !strings.HasPrefix(baseName, "_") {
				if err := CompileSass(sassPath, sassVars, outPath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Palettes returns all theme palettes
func (t *Theme) Palettes() []*Palette {
	// @todo Recomputes on files change
	if len(t.palettes) == 0 {
		// default palette
		t.palettes = append(t.palettes, t.defaultPalette())

		// all palettes defined in theme conf file
		t.palettes = append(t.palettes, t.confPalettes()...)
	}

	return t.palettes
}

// Palette returns the Palette with given name
func (t *Theme) Palette(name string) *Palette {
	for _, p := range t.Palettes() {
		if p.Name == name {
			return p
		}
	}

	return nil
}

// defaultPalette returns default theme palette
func (t *Theme) defaultPalette() *Palette {
	vars, err := t.SassVars()
	if err != nil {
		return nil
	}

	result := NewPalette(defaultPalette)

	for name, val := range vars {
		result.Vars[name] = val
	}

	return result
}

// confPalettes returns all palettes defined in theme conf file
func (t *Theme) confPalettes() []*Palette {
	return t.Conf.Palettes
}

func (t *Theme) setTemplateDir() {
	if t.TemplatesDir == "" {
		dirPath := path.Join(t.Dir, templatesDir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			// no /templates subdir found
			t.TemplatesDir = t.Dir
		} else {
			t.TemplatesDir = dirPath
		}
	}
}

func (t *Theme) setPartialsDir() {
	t.PartialsDir = path.Join(t.TemplatesDir, partialsDir)
}

func (t *Theme) setSassDir() {
	dir := path.Join(t.Dir, sassDir)

	_, err := ioutil.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		// @todo Handle error
		panic(err)
	}

	if os.IsNotExist(err) {
		// no sass dir
		return
	}

	t.SassDir = dir
}

func (t *Theme) setConf() {
	var result Conf

	confPath := path.Join(t.Dir, confFile)
	if _, err := toml.DecodeFile(confPath, &result); err != nil {
		// @todo Handle error
		panic(err)
	}

	result.ID = t.ID

	t.Conf = &result
}
