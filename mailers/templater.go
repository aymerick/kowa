package mailers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/raymond"
)

// TplKind represents a mail template kind
type TplKind string

const (
	tplHTML = TplKind("html")
	tplText = TplKind("txt")
)

var templater *Templater

// Templater handles mail templating
type Templater struct {
	templatesDir string
	layouts      map[TplKind]*raymond.Template
	templates    map[string]*raymond.Template
}

func init() {
	templater = &Templater{
		layouts:   make(map[TplKind]*raymond.Template),
		templates: make(map[string]*raymond.Template),
	}
}

// SetTemplatesDir sets the mail templates directory
func SetTemplatesDir(dir string) {
	templater.templatesDir = dir

	if err := templater.setupLayouts(); err != nil {
		panic(err)
	}
}

// Generate generates a mail template
func (tpl *Templater) Generate(name string, kind TplKind, mailer Mailer) (string, error) {
	template, errT := tpl.getTemplate(name, kind)
	if errT != nil {
		return "", errT
	}

	result, err := template.Exec(mailer)
	if err != nil {
		return "", err
	}

	return result, nil
}

// layout returns a layout
func (tpl *Templater) layout(kind TplKind) *raymond.Template {
	if tpl.layouts[kind] == nil {
		tpl.setupLayouts()
	}

	return tpl.layouts[kind]
}

// getTemplate returns a new template instance
func (tpl *Templater) getTemplate(name string, kind TplKind) (*raymond.Template, error) {
	tplKey := fmt.Sprintf("%s:%s", name, kind)

	if tpl.templates[tplKey] == nil {
		// clone layout
		template := tpl.layout(kind).Clone()

		// adds content partial
		content, err := tpl.templateContent(name, kind)
		if err != nil {
			return nil, err
		}

		template.RegisterPartial("content", content)

		tpl.templates[tplKey] = template
	}

	return tpl.templates[tplKey], nil
}

// templateContent fetches template content
func (tpl *Templater) templateContent(name string, kind TplKind) (string, error) {
	var err error
	var data []byte

	if tpl.templatesDir != "" {
		// fetch from file system
		filePath := path.Join(tpl.templatesDir, fmt.Sprintf("%s.%s.hbs", name, kind))

		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			return "", err
		}
	} else {
		// fetch from embeded assets
		assetPath := fmt.Sprintf("mailers/templates/%s.%s.hbs", name, kind)

		data, err = core.Asset(assetPath)
		if err != nil {
			return "", err
		} else if len(data) == 0 {
			return "", errors.New("Mailer template not found in assets: " + assetPath)
		}
	}

	return string(data), nil
}

// setupLayouts setups layouts
func (tpl *Templater) setupLayouts() error {
	for _, kind := range []TplKind{tplHTML, tplText} {
		if err := tpl.setupLayout(kind); err != nil {
			return err
		}
	}

	return nil
}

func (tpl *Templater) setupLayout(kind TplKind) error {
	content, err := tpl.templateContent("layout", kind)
	if err != nil {
		return err
	}

	layout, err := raymond.Parse(content)
	if err != nil {
		return err
	}

	tpl.layouts[kind] = layout

	return nil
}
