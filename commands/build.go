package commands

import (
	"log"

	"github.com/aymerick/kowa/builder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_THEME      = "minimal"
	DEFAULT_OUTPUT_DIR = "_site"
	DEFAULT_UGLY_URL   = false
)

var buildCmd = &cobra.Command{
	Use:   "build [site_id]",
	Short: "Site builder",
	Long:  `Build a static site.`,
	Run:   buildSite,
}

func initBuilderConf() {
	buildCmd.Flags().StringP("theme", "t", DEFAULT_THEME, "Theme to use")
	viper.BindPFlag("theme", buildCmd.Flags().Lookup("theme"))

	buildCmd.Flags().StringP("output_dir", "o", DEFAULT_OUTPUT_DIR, "Output directory")
	viper.BindPFlag("output_dir", buildCmd.Flags().Lookup("output_dir"))

	buildCmd.Flags().BoolP("urgly_url", "g", DEFAULT_UGLY_URL, "Generate ugly URLs")
	viper.BindPFlag("urgly_url", buildCmd.Flags().Lookup("urgly_url"))
}

func buildSite(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Fatalln("No site id argument provided")
	}

	site := builder.NewSite(args[0])

	log.Printf("Building site '%s' with theme '%s' into %s", args[0], site.Theme, site.GenDir())

	site.Build()
}
