package commands

import (
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/server"
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
	// @todo Check that we are NOT in production

	// Insert users
	userJeanClaude := models.User{
		Id:        bson.NewObjectId(),
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		CreatedAt: time.Now(),
	}
	models.UsersCol().Insert(&userJeanClaude)

	userHenry := models.User{
		Id:        bson.NewObjectId(),
		FirstName: "Henry",
		LastName:  "Kanan",
		CreatedAt: time.Now(),
	}
	models.UsersCol().Insert(&userHenry)

	// Insert sites
	var site models.Site

	site = models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userJeanClaude.Id,
		CreatedAt:   time.Now(),
		Name:        "My site",
		Tagline:     "So powerfull !",
		Description: "You will be astonished by what my site is about",
	}
	models.SitesCol().Insert(&site)

	site = models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userJeanClaude.Id,
		CreatedAt:   time.Now(),
		Name:        "My second site",
		Tagline:     "Very interesting",
		Description: "Our projects are so importants, please help us",
	}
	models.SitesCol().Insert(&site)

	site = models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userHenry.Id,
		CreatedAt:   time.Now(),
		Name:        "Ultimate petanque",
		Tagline:     "La petanque comme vous ne l'avez jamais vu",
		Description: "C'est vraiment le sport du futur. Messieurs, preparez vos boules !",
	}
	models.SitesCol().Insert(&site)

	// Insert oauth client
	oauthStorage := server.NewOAuthStorage()
	oauthStorage.SetupDefaultClient()
}
