package themes

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/wellington/go-libsass"
)

// absolute path to themes directory
var dir string

// Dir returns path to themes directory
func Dir() string {
	if dir != "" {
		return dir
	}

	return viper.GetString("themes_dir")
}

// SetDir sets path to themes directory
func SetDir(path string) {
	dir = path
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
