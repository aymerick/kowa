package commands

import (
	"github.com/aymerick/kowa/builder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_THEME      = "minimal"
	DEFAULT_OUTPUT_DIR = "_site"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Site builder",
	Long:  `Build a static site.`,
	Run:   buildSite,
}

func initBuilderConf() {
	buildCmd.Flags().StringP("theme", "t", DEFAULT_THEME, "Theme to use")
	viper.BindPFlag("theme", buildCmd.Flags().Lookup("theme"))

	buildCmd.Flags().StringP("output_dir", "o", DEFAULT_OUTPUT_DIR, "Output directory")
	viper.BindPFlag("output_dir", buildCmd.Flags().Lookup("output_dir"))
}

func buildSite(cmd *cobra.Command, args []string) {
	site := builder.NewSite()
	site.Build()
}
