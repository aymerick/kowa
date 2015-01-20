package builder

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"
)

type PageMeta struct {
	Title       string
	Description string
}

type Page struct {
	Meta      *PageMeta
	BodyClass string
	Head      string
	Footer    string
	Content   interface{}
}

func Build() {
	layout := template.Must(template.ParseFiles(templatePath("layout")))

	activitiesPageTpl := pageTemplate("activities", layout)
	contactPageTpl := pageTemplate("contact", layout)

	_ = activitiesPageTpl.Execute(os.Stdout, Page{
		Meta: &PageMeta{
			Title:       "Activities",
			Description: "Activities test page",
		},
		BodyClass: "activities",
		Content:   []string{"one", "two", "three<br />", "four"},
	})

	_ = contactPageTpl.Execute(os.Stdout, Page{
		Meta: &PageMeta{
			Title:       "Contact",
			Description: "Contact test page",
		},
		BodyClass: "contact",
		Content:   "Contact<br /> \\o/",
	})
}

func pageTemplate(name string, layout *template.Template) *template.Template {
	result := template.Must(layout.Clone())
	b, _ := ioutil.ReadFile(templatePath(name))
	_, _ = result.New("content").Parse(string(b))

	return result
}

func templatePath(tplName string) string {
	return path.Join(viper.GetString("working_dir"), "themes", viper.GetString("theme"), fmt.Sprintf("%s.html", tplName))
}
