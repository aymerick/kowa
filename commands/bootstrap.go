package commands

import (
	"fmt"
	"math/rand"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

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

	rand.Seed(8941)

	now := time.Now()
	lastWeek := now.Add(-24 * 7 * time.Hour)

	// Insert oauth client
	oauthStorage := server.NewOAuthStorage()
	oauthStorage.SetupDefaultClient()

	db := models.NewDBSession()

	password, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		panic("Arg")
	}

	// Insert users
	userJeanClaude := models.User{
		Id:        "test",
		CreatedAt: lastWeek,
		Email:     "trucmush@wanadoo.fr",
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		Password:  string(password),
	}
	db.UsersCol().Insert(&userJeanClaude)

	userHenry := models.User{
		Id:        "hkanan",
		CreatedAt: lastWeek,
		Email:     "henrykanan@yahoo.com",
		FirstName: "Henry",
		LastName:  "Kanan",
		Password:  string(password),
	}
	db.UsersCol().Insert(&userHenry)

	// Insert sites
	siteJC1 := models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userJeanClaude.Id,
		CreatedAt:   lastWeek,
		Name:        "My site",
		Tagline:     "So powerfull !",
		Description: "You will be astonished by what my site is about",
	}
	db.SitesCol().Insert(&siteJC1)

	siteJC2 := models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userJeanClaude.Id,
		CreatedAt:   lastWeek,
		Name:        "My second site",
		Tagline:     "Very interesting",
		Description: "Our projects are so importants, please help us",
	}
	db.SitesCol().Insert(&siteJC2)

	siteH := models.Site{
		Id:          bson.NewObjectId(),
		UserId:      userHenry.Id,
		CreatedAt:   lastWeek,
		Name:        "Ultimate petanque",
		Tagline:     "La petanque comme vous ne l'avez jamais vu",
		Description: "C'est vraiment le sport du futur. Messieurs, preparez vos boules !",
	}
	db.SitesCol().Insert(&siteH)

	// Insert posts
	var post models.Post

	for i := 1; i <= 30; i++ {
		nbDur := time.Duration(i)

		post = models.Post{
			Id:          bson.NewObjectId(),
			CreatedAt:   lastWeek.Add(time.Hour * nbDur),
			UpdatedAt:   lastWeek.Add(time.Hour*nbDur + 30),
			SiteId:      siteJC1.Id,
			PublishedAt: lastWeek.Add(time.Hour*nbDur + 30),
			Title:       fmt.Sprintf("Post %d", i),
			Body:        fmt.Sprintf(MD_FIXTURES[rand.Intn(len(MD_FIXTURES))]),
		}
		db.PostsCol().Insert(&post)
	}

	post = models.Post{
		Id:          bson.NewObjectId(),
		CreatedAt:   lastWeek.Add(time.Hour),
		UpdatedAt:   lastWeek.Add(time.Hour + 30),
		SiteId:      siteJC2.Id,
		PublishedAt: lastWeek.Add(time.Hour + 30),
		Title:       "This is a lonely",
		Body:        "It appears on my second website.",
	}
	db.PostsCol().Insert(&post)

	post = models.Post{
		Id:          bson.NewObjectId(),
		CreatedAt:   lastWeek.Add(48 * time.Hour),
		UpdatedAt:   lastWeek.Add(48*time.Hour + 30),
		SiteId:      siteH.Id,
		PublishedAt: lastWeek.Add(48*time.Hour + 30),
		Title:       "Hi, I am Henry",
		Body:        "Je me présente, je m'appelle Henry. Je voudrais bien réussir ma vie, être aimé. Être beau, gagner de l'argent. Puis surtout être intelligent. Mais pour tout ça il faudrait que je bosse à plein temps",
	}
	db.PostsCol().Insert(&post)

	// Insert events

	// Insert pages

	// Insert actions
}
