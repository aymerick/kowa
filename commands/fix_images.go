package commands

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/models"
)

var fixImagesCmd = &cobra.Command{
	Use:   "fix_images [site_id]",
	Short: "Fix images",
	Long:  `Fix all images for given site.`,
	Run:   fixImages,
}

func fixImages(cmd *cobra.Command, args []string) {
	dbSession := models.NewDBSession()

	if len(args) < 1 {
		cmd.Usage()
		log.Fatalln("ERROR: No site id argument provided")
	}

	site := dbSession.FindSite(args[0])
	if site == nil {
		cmd.Usage()
		log.Fatalln("ERROR: Site not found:" + args[0])
	}

	for _, image := range *site.FindAllImages() {
		prefix := "/upload/" + site.ID + "/"

		if strings.HasPrefix(image.Path, prefix) {
			fixedPath := strings.TrimPrefix(image.Path, prefix)

			log.Printf("Fixing path: %s => %s", image.Path, fixedPath)

			dbSession.ImagesCol().UpdateId(image.ID, bson.M{"$set": bson.M{"path": fixedPath}})
		}
	}
}
