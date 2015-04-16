package mailers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/aymerick/kowa/core"
)

type TplKind string

const (
	TPL_HTML = TplKind("html")
	TPL_TEXT = TplKind("txt")
)

var templater *Templater

type Templater struct {
	templatesDir string
	layouts      map[TplKind]*template.Template
	templates    map[string]string
}

func init() {
	templater = &Templater{
		layouts:   make(map[TplKind]*template.Template),
		templates: make(map[string]string),
	}
}

func SetTemplatesDir(dir string) {
	templater.templatesDir = dir

	if err := templater.setupLayouts(); err != nil {
		panic(err)
	}
}

// Generates template
func (tpl *Templater) Generate(name string, kind TplKind, mailer Mailer) (string, error) {
	var result bytes.Buffer

	// get template instance
	tplInstance := template.Must(template.Must(tpl.getTemplate(name, kind, mailer)).Clone())

	// execute template
	if err := tplInstance.Execute(&result, mailer); err != nil {
		return "", err
	} else {
		return result.String(), nil
	}
}

// Get a layout
func (tpl *Templater) layout(name string, kind TplKind) *template.Template {
	if tpl.layouts[kind] == nil {
		tpl.setupLayouts()
	}

	return tpl.layouts[kind]
}

// Returns a new template instance
func (tpl *Templater) getTemplate(name string, kind TplKind, mailer Mailer) (*template.Template, error) {
	// clone layout
	result, errL := tpl.layout(name, kind).Clone()
	if errL != nil {
		return nil, errL
	}

	// @todo setup FuncMap
	// result.Funcs(mailer.FuncMap())

	// parse template
	_, err := result.New("content").Parse(tpl.templateContent(name, kind))
	if err != nil {
		return nil, err
	}

	return result, err
}

// Fetch template content
func (tpl *Templater) templateContent(name string, kind TplKind) string {
	tplKey := fmt.Sprintf("%s:%s", name, kind)

	if tpl.templates[tplKey] == "" {
		var err error
		var data []byte

		if tpl.templatesDir != "" {
			// fetch from file system
			filePath := path.Join(tpl.templatesDir, fmt.Sprintf("%s.%s", name, kind))

			data, err = ioutil.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			// fetch from embeded assets
			assetPath := fmt.Sprintf("mailers/templates/%s.%s", name, kind)

			data, err = core.Asset(assetPath)
			if err != nil {
				panic(err)
			} else if len(data) == 0 {
				panic("Mailer template not found in assets: " + assetPath)
			}
		}

		tpl.templates[tplKey] = string(data)
	}

	return tpl.templates[tplKey]
}

// Returns template file path
// func (tpl *Templater) templatePath(name string, kind TplKind) string {
// 	return path.Join(tpl.templatesDir, fmt.Sprintf("%s.%s", name, kind))
// }

// Setup layouts
func (tpl *Templater) setupLayouts() error {
	// fetch html layout
	htmlLayout, err := template.New("layout").Parse(tpl.templateContent("layout", TPL_HTML))
	if err != nil {
		return err
	}

	tpl.layouts[TPL_HTML] = htmlLayout

	// fetch text layout
	textLayout, errT := template.New("layout").Parse(tpl.templateContent("layout", TPL_TEXT))
	if errT != nil {
		return errT
	}

	tpl.layouts[TPL_TEXT] = textLayout

	// @todo load partials

	return nil
}
