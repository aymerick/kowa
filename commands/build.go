package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/builder"
	"github.com/aymerick/kowa/models"
)

const (
	DEFAULT_UGLY_URL = false
)

var buildCmd = &cobra.Command{
	Use:   "build [site_id]",
	Short: "Site builder",
	Long:  `Build a static site.`,
	Run:   buildSite,
}

func initBuilderConf() {
	buildCmd.Flags().StringP("theme", "t", builder.DEFAULT_THEME, "Theme to use")
	viper.BindPFlag("theme", buildCmd.Flags().Lookup("theme"))

	buildCmd.Flags().BoolP("ugly_url", "g", DEFAULT_UGLY_URL, "Generate ugly URLs")
	viper.BindPFlag("ugly_url", buildCmd.Flags().Lookup("ugly_url"))
}

func buildSite(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Fatalln("No site id argument provided")
	}

	// get site
	site := models.NewDBSession().FindSite(args[0])
	if site == nil {
		cmd.Usage()
		log.Fatalln("Site not found:" + args[0])
	}

	// builder config
	config := &builder.SiteBuilderConfig{
		WorkingDir: viper.GetString("working_dir"),
		OutputDir:  viper.GetString("output_dir"),
		Theme:      viper.GetString("theme"),
		UglyURL:    viper.GetBool("ugly_url"),
	}

	siteBuilder := builder.NewSiteBuilder(site, config)

	log.Printf("Building site '%s' with theme '%s' into %s", args[0], siteBuilder.Theme(), siteBuilder.GenDir())

	// build site
	if siteBuilder.Build(); siteBuilder.HaveError() {
		siteBuilder.DumpErrors()
		siteBuilder.DumpLayout()
	} else {
		if viper.GetBool("serve_output") {
			// server site
			serve(siteBuilder, viper.GetInt("serve_output_port"))
		}
	}
}

func serve(siteBuilder *builder.SiteBuilder, port int) {
	log.Printf("Serving built site from: " + siteBuilder.GenDir())

	log.Printf("Web Server is available at http://127.0.0.1:%d\n", port)
	log.Printf("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(siteBuilder.GenDir()))))
}
