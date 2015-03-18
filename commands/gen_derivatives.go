package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/models"
)

var genDerivativesCmd = &cobra.Command{
	Use:   "gen_derivatives [site_id]",
	Short: "Generate derivatives",
	Long:  `Regenerate all images derivatives for given site.`,
	Run:   genDerivatives,
}

func genDerivatives(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Fatalln("ERROR: No site id argument provided")
	}

	if viper.GetString("upload_dir") == "" {
		cmd.Usage()
		log.Fatalln("ERROR: The upload_dir setting is mandatory")
	}

	// get site
	site := models.NewDBSession().FindSite(args[0])
	if site == nil {
		cmd.Usage()
		log.Fatalln("ERROR: Site not found:" + args[0])
	}

	// generate derivatives
	for _, image := range *site.FindAllImages() {
		if err := image.GenerateDerivatives(true); err != nil {
			log.Printf("ERROR: Failed to generate image: %v", err)
		}
	}
}
