package builder

import (
	"fmt"
	"path"

	"github.com/spf13/viper"
)

func templatePath(tplName string) string {
	return path.Join(viper.GetString("working_dir"), "themes", viper.GetString("theme"), fmt.Sprintf("%s.html", tplName))
}

func partialsPath() string {
	return path.Join(viper.GetString("working_dir"), "themes", viper.GetString("theme"), "partials")
}
