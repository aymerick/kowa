package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aymerick/kowa/builder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_THEME      = "minimal"
	DEFAULT_OUTPUT_DIR = "_site"
	DEFAULT_UGLY_URL   = false
	DEFAULT_SERVE      = false
	DEFAULT_SERVE_PORT = 48910
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

	buildCmd.Flags().BoolP("ugly_url", "g", DEFAULT_UGLY_URL, "Generate ugly URLs")
	viper.BindPFlag("ugly_url", buildCmd.Flags().Lookup("ugly_url"))

	buildCmd.Flags().BoolP("serve", "s", DEFAULT_SERVE, "Start a server to test built site")
	viper.BindPFlag("serve", buildCmd.Flags().Lookup("serve"))

	serverCmd.Flags().IntP("serve_port", "t", DEFAULT_SERVE_PORT, "Port to test built site")
	viper.BindPFlag("serve_port", serverCmd.Flags().Lookup("serve_port"))
}

func buildSite(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Fatalln("No site id argument provided")
	}

	siteBuilder := builder.NewSiteBuilder(args[0])

	log.Printf("Building site '%s' with theme '%s' into %s", args[0], viper.GetString("theme"), siteBuilder.GenDir())

	// build site
	siteBuilder.Build()

	if viper.GetBool("serve") {
		// server site
		serve(siteBuilder, viper.GetInt("serve_port"))
	}
}

func serve(siteBuilder *builder.SiteBuilder, port int) {
	log.Printf("Serving built site from: " + siteBuilder.GenDir())

	log.Printf("Web Server is available at http://127.0.0.1:%d\n", port)
	log.Printf("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(siteBuilder.GenDir()))))
}
