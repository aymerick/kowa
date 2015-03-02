package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/aymerick/kowa/builder"
	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
)

var addSiteCmd = &cobra.Command{
	Use:   "add_site [id] [user_id]",
	Short: "Add a new site",
	Long:  `Insert a new site in database.`,
	Run:   addSite,
}

func addSite(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		log.Fatalln("Missing arguments")
	}

	dbSession := models.NewDBSession()

	if site := dbSession.FindSite(args[0]); site != nil {
		log.Fatalln("There is already a site with that id: " + args[0])
	}

	site := &models.Site{
		Id:     args[0],
		UserId: args[1],
		Theme:  builder.DEFAULT_THEME,
	}

	if err := dbSession.CreateSite(site); err != nil {
		log.Fatalln(fmt.Sprintf("Failed to create site: %v", err))
	}

	utils.AppEnsureSiteUploadDir(site.Id)
}
