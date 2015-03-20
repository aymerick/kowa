package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/models"
)

var addSiteCmd = &cobra.Command{
	Use:   "add_site [id] [site name] [user id]",
	Short: "Add a new site",
	Long:  `Insert a new site in database.`,
	Run:   addSite,
}

func addSite(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		cmd.Usage()
		log.Fatalln("Missing arguments")
	}

	dbSession := models.NewDBSession()

	if site := dbSession.FindSite(args[0]); site != nil {
		log.Fatalln("There is already a site with that id: " + args[0])
	}

	// @todo FIXME !
	baseUrl := fmt.Sprintf("%s:%d/%s", core.DEFAULT_BASEURL, viper.GetInt("serve_output_port"), args[0])

	site := &models.Site{
		Id:           args[0],
		Name:         args[1],
		UserId:       args[2],
		Lang:         core.DEFAULT_LANG,
		Theme:        core.DEFAULT_THEME,
		BaseUrl:      baseUrl,
		NameInNavBar: true,
	}

	if err := dbSession.CreateSite(site); err != nil {
		log.Fatalln(fmt.Sprintf("Failed to create site: %v", err))
	}

	core.EnsureSiteUploadDir(site.Id)

	// build site
	buildSite(site)
}
