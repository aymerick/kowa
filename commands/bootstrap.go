package commands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/server"
	"github.com/aymerick/kowa/utils"
)

const (
	IMAGE_FIXTURES_DIR = "/fixtures"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap Kowa",
	Long:  `Bootstrap database with fake data`,
	Run:   bootstrap,
}

func bootstrap(cmd *cobra.Command, args []string) {
	// @todo Check that we are NOT in production

	utils.AppEnsureUploadDir()

	rand.Seed(time.Now().UTC().UnixNano())

	now := time.Now()
	lastMonth := now.Add(-31 * 24 * time.Hour)

	//
	// Insert oauth client
	//

	oauthStorage := server.NewOAuthStorage()
	oauthStorage.SetupDefaultClient()

	db := models.NewDBSession()

	//
	// Insert users
	//

	password, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		panic("Arg")
	}

	userJeanClaude := models.User{
		Id:        "test",
		CreatedAt: lastMonth,
		Email:     "trucmush@wanadoo.fr",
		FirstName: "Jean-Claude",
		LastName:  "Trucmush",
		Password:  string(password),
	}
	db.UsersCol().Insert(&userJeanClaude)

	userHenry := models.User{
		Id:        "hkanan",
		CreatedAt: lastMonth,
		Email:     "henrykanan@yahoo.com",
		FirstName: "Henry",
		LastName:  "Kanan",
		Password:  string(password),
	}
	db.UsersCol().Insert(&userHenry)

	//
	// Insert sites
	//

	siteJC1 := models.Site{
		Id:          "site_1",
		UserId:      userJeanClaude.Id,
		CreatedAt:   lastMonth,
		Name:        "My site",
		Tagline:     "So powerfull !",
		Description: "You will be astonished by what my site is about",
	}
	if err := db.SitesCol().Insert(&siteJC1); err != nil {
		panic(err)
	}

	utils.AppEnsureSiteUploadDir(siteJC1.Id)

	siteJC2 := models.Site{
		Id:          "site_2",
		UserId:      userJeanClaude.Id,
		CreatedAt:   lastMonth,
		Name:        "My second site",
		Tagline:     "Very interesting",
		Description: "Our projects are so importants, please help us",
	}
	db.SitesCol().Insert(&siteJC2)

	utils.AppEnsureSiteUploadDir(siteJC2.Id)

	siteH := models.Site{
		Id:          "ultimate",
		UserId:      userHenry.Id,
		CreatedAt:   lastMonth,
		Name:        "Ultimate petanque",
		Tagline:     "La petanque comme vous ne l'avez jamais vu",
		Description: "C'est vraiment le sport du futur. Messieurs, preparez vos boules !",
	}
	db.SitesCol().Insert(&siteH)

	utils.AppEnsureSiteUploadDir(siteH.Id)

	sites := []models.Site{siteJC1, siteJC2, siteH}

	//
	// Insert images
	//

	siteJC1Images := []models.Image{}

	currentDir, errWd := os.Getwd()
	if errWd != nil {
		panic(errWd)
	}

	imgFiles, errDir := ioutil.ReadDir(path.Join(currentDir, "/client/public", IMAGE_FIXTURES_DIR))
	if errDir != nil {
		panic(errDir)
	}

	for i, imgFile := range imgFiles {
		if !imgFile.IsDir() && !models.IsDerivativePath(imgFile.Name()) {
			fileExt := path.Ext(imgFile.Name())
			switch fileExt {
			case ".png", ".jpg", ".gif", ".PNG", ".JPG", ".GIF":
				var fileType string

				switch fileExt {
				case ".png", ".PNG":
					fileType = "image/png"
				case ".jpg", ".JPG":
					fileType = "image/jpeg"
				case ".gif", ".GIF":
					fileType = "image/gif"
				}

				// insert image in all sites
				for j, site := range sites {
					nbHours := time.Duration(i)

					img := models.Image{
						Id:        bson.NewObjectId(),
						CreatedAt: lastMonth.Add(time.Hour * nbHours),
						UpdatedAt: lastMonth.Add(time.Hour * nbHours),
						SiteId:    site.Id,
						Path:      path.Join(IMAGE_FIXTURES_DIR, imgFile.Name()),
						Name:      imgFile.Name(),
						Size:      imgFile.Size(),
						Type:      fileType,
					}
					db.ImagesCol().Insert(&img)

					if site.Id == siteJC1.Id {
						siteJC1Images = append(siteJC1Images, img)
					}

					if j == 0 {
						errThumb := img.GenerateDerivatives()
						if errThumb != nil {
							panic(errThumb)
						}
					}
				}
			}
		}
	}

	//
	// Insert posts
	//

	var post models.Post

	for i := 1; i <= 30; i++ {
		nbDays := time.Duration(i)

		pubDate := lastMonth.Add(time.Hour*24*nbDays + 30)

		post = models.Post{
			Id:          bson.NewObjectId(),
			CreatedAt:   lastMonth.Add(time.Hour * 24 * nbDays),
			UpdatedAt:   lastMonth.Add(time.Hour*24*nbDays + 30),
			SiteId:      siteJC1.Id,
			PublishedAt: &pubDate,
			Title:       fmt.Sprintf("Post %d", i),
			Body:        fmt.Sprintf(MD_FIXTURES[rand.Intn(len(MD_FIXTURES))]),
			Cover:       siteJC1Images[rand.Intn(len(siteJC1Images))].Id,
		}
		db.PostsCol().Insert(&post)
	}

	pubDate := lastMonth.Add(time.Hour + 30)

	post = models.Post{
		Id:          bson.NewObjectId(),
		CreatedAt:   lastMonth.Add(time.Hour),
		UpdatedAt:   lastMonth.Add(time.Hour + 30),
		SiteId:      siteJC2.Id,
		PublishedAt: &pubDate,
		Title:       "This is a lonely post",
		Body:        "It appears on my second website.",
	}
	db.PostsCol().Insert(&post)

	post = models.Post{
		Id:          bson.NewObjectId(),
		CreatedAt:   lastMonth.Add(48 * time.Hour),
		UpdatedAt:   lastMonth.Add(48*time.Hour + 30),
		SiteId:      siteH.Id,
		PublishedAt: &pubDate,
		Title:       "Hi, I am Henry",
		Body:        "Je me présente, je m'appelle Henry. Je voudrais bien réussir ma vie, être aimé. Être beau, gagner de l'argent. Puis surtout être intelligent. Mais pour tout ça il faudrait que je bosse à plein temps",
	}
	db.PostsCol().Insert(&post)

	// @todo Insert events

	//
	// Insert pages
	//

	var page models.Page

	for i := 1; i <= 30; i++ {
		nbDays := time.Duration(i)

		page = models.Page{
			Id:        bson.NewObjectId(),
			CreatedAt: lastMonth.Add(time.Hour * 24 * nbDays),
			UpdatedAt: lastMonth.Add(time.Hour*24*nbDays + 30),
			SiteId:    siteJC1.Id,
			Title:     fmt.Sprintf("Page %d", i),
			Tagline:   fmt.Sprintf("This a fantastic tagline %d", i),
			Body:      fmt.Sprintf(MD_FIXTURES[rand.Intn(len(MD_FIXTURES))]),
			Cover:     siteJC1Images[rand.Intn(len(siteJC1Images))].Id,
		}
		db.PagesCol().Insert(&page)
	}

	page = models.Page{
		Id:        bson.NewObjectId(),
		CreatedAt: lastMonth.Add(time.Hour),
		UpdatedAt: lastMonth.Add(time.Hour + 30),
		SiteId:    siteJC2.Id,
		Title:     "This is a lonely page",
		Tagline:   "That page rocks so much",
		Body:      "It appears on my second website.",
	}
	db.PagesCol().Insert(&page)

	//
	// Insert activities
	//

	var activity models.Activity

	for i := 1; i <= 30; i++ {
		nbDays := time.Duration(i)

		activity = models.Activity{
			Id:        bson.NewObjectId(),
			CreatedAt: lastMonth.Add(time.Hour * 24 * nbDays),
			UpdatedAt: lastMonth.Add(time.Hour*24*nbDays + 30),
			SiteId:    siteJC1.Id,
			Title:     fmt.Sprintf("Activity %d", i),
			Body:      fmt.Sprintf(MD_FIXTURES[rand.Intn(len(MD_FIXTURES))]),
			Cover:     siteJC1Images[rand.Intn(len(siteJC1Images))].Id,
		}
		db.ActivitiesCol().Insert(&activity)
	}
}
