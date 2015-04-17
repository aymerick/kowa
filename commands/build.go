package commands

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/builder"
	"github.com/aymerick/kowa/models"
)

var buildCmd = &cobra.Command{
	Use:   "build [site_id]",
	Short: "Build site",
	Long:  `Build a static site.`,
	Run:   buildSiteCmd,
}

func buildSiteCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Fatalln("ERROR: No site id argument provided")
	}

	checkAndOutputsGlobalFlags()

	// get site
	site := models.NewDBSession().FindSite(args[0])
	if site == nil {
		cmd.Usage()
		log.Fatalln("ERROR: Site not found:" + args[0])
	}

	// build site
	siteBuilder := buildSite(site)
	if siteBuilder.HaveError() {
		siteBuilder.DumpErrors()
		siteBuilder.DumpLayout()
	} else {
		if viper.GetBool("serve_output") {
			// server site
			serve(siteBuilder, viper.GetInt("serve_output_port"))
		}
	}
}

func buildSite(site *models.Site) *builder.SiteBuilder {
	// builder config
	config := &builder.SiteBuilderConfig{
		ThemesDir: viper.GetString("themes_dir"),
		OutputDir: path.Join(viper.GetString("output_dir"), site.Id),
	}

	siteBuilder := builder.NewSiteBuilder(site, config)

	log.Printf("Building site '%s' with theme '%s' into %s", site.Id, site.Theme, config.OutputDir)

	startTime := time.Now()

	// build
	if siteBuilder.Build(); siteBuilder.HaveError() {
		log.Println("Failed to build site")
	} else {
		// update BuiltAt anchor
		site.SetBuiltAt(time.Now())

		log.Printf("Site build in %v ms\n", int(1000*time.Since(startTime).Seconds()))
	}

	return siteBuilder
}

func serve(siteBuilder *builder.SiteBuilder, port int) {
	servePath, _ := path.Split(siteBuilder.Config().OutputDir)

	log.Printf("Serving built site from: " + servePath)

	log.Printf("Web Server is available at http://127.0.0.1:%d\n", port)
	log.Printf("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(servePath))))
}
