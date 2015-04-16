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

// Returns a new template instance
func (tpl *Templater) getTemplate(name string, kind TplKind, mailer Mailer) (*template.Template, error) {
	tplKey := fmt.Sprintf("%s:%s", name, kind)

	if tpl.templates[tplKey] == "" {
		// read template file
		binData, err := ioutil.ReadFile(tpl.templatePath(name, kind))
		if err != nil {
			return nil, err
		}

		tpl.templates[tplKey] = string(binData)
	}

	// clone layout
	result, errL := tpl.layouts[kind].Clone()
	if errL != nil {
		return nil, errL
	}

	// @todo setup FuncMap
	// result.Funcs(mailer.FuncMap())

	// parse template
	_, err := result.New("content").Parse(tpl.templates[tplKey])
	if err != nil {
		return nil, err
	}

	return result, err
}

// Returns template file path
func (tpl *Templater) templatePath(name string, kind TplKind) string {
	return path.Join(tpl.templatesDir, fmt.Sprintf("%s.%s", name, kind))
}

// Setup layouts
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

	// @todo load partials

	return nil
}
