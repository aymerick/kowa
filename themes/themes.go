package themes

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/aymerick/kowa/helpers"
	"github.com/spf13/viper"
	"github.com/wellington/go-libsass"
)

// absolute path to themes directory
var dir string

// all installed themes
var themes = map[string]*Theme{}

// Dir returns path to themes directory
func Dir() string {
	if dir != "" {
		return dir
	}

	return viper.GetString("themes_dir")
}

// SetDir sets path to themes directory
func SetDir(path string) {
	// @todo Reset themes var
	dir = path
}

// Get returns theme instance for given theme id
func Get(id string) *Theme {
	return themes[id]
}

// All returns all installed themes
func All() map[string]*Theme {
	// @todo Recompute on file change
	if len(themes) == 0 {
		files, err := ioutil.ReadDir(Dir())
		if err != nil {
			// @todo Handle error
			panic(err)
		}

		for _, file := range files {
			filePath := path.Join(Dir(), file.Name())

			if file.Mode()&os.ModeSymlink != 0 {
				// ioutil.ReadDir() performs os.Lstat() calls that do not follow symlinks... so fix that
				file, err = os.Stat(filePath)
				if err != nil {
					// @todo Handle error
					panic(err)
				}
			}

			// check that this is a directory with a theme.toml conf file
			if file.IsDir() {
				conf := path.Join(filePath, confFile)
				if _, err := os.Stat(conf); !os.IsNotExist(err) {

					themeID := helpers.FileBase(file.Name())

					themes[themeID] = New(themeID)
				}
			}
		}
	}

	return themes
}

// Exist returns true if a theme with given id exists
func Exist(id string) bool {
	return All()[id] != nil
}

// AllConf returns all installed themes configurations
func AllConf() []*Conf {
	var result []*Conf

	for _, t := range All() {
		result = append(result, t.Conf)
	}

	return result
}

// CompileSass compiles given sass file
func CompileSass(sassFilePath string, sassVars string, outPath string) error {
	ctx := libsass.Context{
		BuildDir:     filepath.Dir(outPath),
		MainFile:     sassFilePath,
		IncludePaths: []string{filepath.Dir(sassFilePath)},
	}

	ctx.Imports.Init()

	if sassVars != "" {
		// overwrite _variables.scss partial with given sass code
		ctx.Imports.Add("", "variables", []byte(sassVars))
	}

	// create directory
	dirPath, _ := path.Split(outPath)
	if err := os.MkdirAll(dirPath, 0755); (err != nil) && !os.IsExist(err) {
		log.Printf("Failed to create dir: '%s'", dirPath)
		return err
	}

	// create output file
	out, err := os.Create(outPath)
	if err != nil {
		log.Printf("Failed to create file: '%s'", outPath)
		return err
	}
	defer out.Close()

	// open sass file
	sassFile, err := os.Open(sassFilePath)
	if err != nil {
		return err
	}
	defer sassFile.Close()

	// compile to CSS
	if err := ctx.Compile(sassFile, out); err != nil {
		log.Printf("Failed to compile sass file: '%s'", sassFilePath)
		return err
	}

	return nil
}
