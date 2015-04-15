package mailers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"
)

type TplKind string

const (
	TPL_HTML = TplKind("html")
	TPL_TEXT = TplKind("txt")
)

var templater *Templater

type Templater struct {
	layouts      map[TplKind]*template.Template
	templatesDir string
}

func init() {
	templater = &Templater{
		layouts: make(map[TplKind]*template.Template),
	}
}

func SetTemplatesDir(dir string) {
	templater.templatesDir = dir

	if err := templater.setupLayouts(); err != nil {
		panic(err)
	}
}

func (tpl *Templater) setupLayouts() error {
	if tpl.templatesDir == "" {
		// @todo Embeds assets
		panic("NOT IMPLEMENTED - @todo Use built-in templates")
	}

	// parse layouts
	htmlLayout, err := template.ParseFiles(tpl.templatePath("layout", TPL_HTML))
	if err != nil {
		return err
	}

	tpl.layouts[TPL_HTML] = htmlLayout

	textLayout, err := template.ParseFiles(tpl.templatePath("layout", TPL_TEXT))
	if err != nil {
		return err
	}

	tpl.layouts[TPL_TEXT] = textLayout

	// @todo setup FuncMap
	// @todo load partials

	return nil
}

// Returns template file path
func (tpl *Templater) templatePath(name string, kind TplKind) string {
	return path.Join(tpl.templatesDir, fmt.Sprintf("%s.%s", name, kind))
}

// Returns template
func (tpl *Templater) template(name string, kind TplKind) (*template.Template, error) {
	result := template.Must(tpl.layouts[kind].Clone())

	// add "content" template to main layout
	binData, err := ioutil.ReadFile(tpl.templatePath(name, kind))
	if err == nil {
		_, err = result.New("content").Parse(string(binData))
	}

	return result, err
}

// Generates template
func (tpl *Templater) generate(name string, kind TplKind, data interface{}) (string, error) {
	var result bytes.Buffer

	if err := template.Must(tpl.template(name, kind)).Execute(&result, data); err != nil {
		return "", err
	} else {
		return result.String(), nil
	}
}
