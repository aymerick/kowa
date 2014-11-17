package commands

import (
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap Kowa",
	Long:  `Creates records in database`,
	Run:   bootstrap,
}

func bootstrap(cmd *cobra.Command, args []string) {
	// Insert user
	user := models.User{
		Id:        bson.NewObjectId(),
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		CreatedAt: time.Now(),
	}
	models.UsersCol().Insert(&user)

	// Insert site
	site := models.Site{
		Id:          bson.NewObjectId(),
		UserId:      user.Id,
		CreatedAt:   time.Now(),
		Name:        "My site",
		Tagline:     "So powerfull !",
		Description: "You will be astonished by what my site is about",
	}
	models.SitesCol().Insert(&site)
}
