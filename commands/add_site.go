package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

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

	site := &models.Site{
		ID:           args[0],
		Name:         args[1],
		UserID:       args[2],
		Lang:         core.DefaultLang,
		Theme:        core.DefaultTheme,
		NameInNavBar: true,
	}

	if domain := core.DefaultDomain(); domain != "" {
		site.Domain = domain
	} else {
		site.CustomURL = core.BaseUrl(args[0])
	}

	if err := dbSession.CreateSite(site); err != nil {
		log.Fatalln(fmt.Sprintf("Failed to create site: %v", err))
	}

	core.EnsureSiteUploadDir(site.ID)

	// build site
	buildSite(site)
}
