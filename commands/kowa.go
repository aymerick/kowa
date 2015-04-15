package commands

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/helpers"
)

const (
	DEFAULT_MONGODB_URI    = "mongodb://localhost:27017/"
	DEFAULT_MONGODB_DBNAME = "kowa"
	DEFAULT_SERVE          = false
	DEFAULT_SERVE_PORT     = 48910

	DEFAULT_SERVICE_NAME      = "Kowa"
	DEFAULT_SERVICE_URL       = "https://github.com/aymerick/kowa" // @todo FIXME
	DEFAULT_SERVICE_COPYRIGHT = "Copyright @ 2015 Kowa - All rights reserved"

	DEFAULT_OUTPUT_DIR = "_sites"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kowa",
	Short: "Kowa static website manager",
	Long:  "Kowa static website manager",
}

func initKowaConf() {
	cobra.OnInitialize(setupConfig)

	// config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kowa/config.toml)")

	// client app
	rootCmd.PersistentFlags().StringP("upload_dir", "u", "", "Uploaded files are stored in that directory (MANDATORY)")
	viper.BindPFlag("upload_dir", rootCmd.PersistentFlags().Lookup("upload_dir"))

	// mongodb database
	rootCmd.PersistentFlags().StringP("mongodb_uri", "d", DEFAULT_MONGODB_URI, "Uri to connect to mongoDB")
	viper.BindPFlag("mongodb_uri", rootCmd.PersistentFlags().Lookup("mongodb_uri"))

	rootCmd.PersistentFlags().StringP("mongodb_dbname", "n", DEFAULT_MONGODB_DBNAME, "MongoDB database name")
	viper.BindPFlag("mongodb_dbname", rootCmd.PersistentFlags().Lookup("mongodb_dbname"))

	// builder
	rootCmd.PersistentFlags().StringP("themes_dir", "t", "", "Themes directory (MANDATORY)")
	viper.BindPFlag("themes_dir", rootCmd.PersistentFlags().Lookup("themes_dir"))

	rootCmd.PersistentFlags().StringP("output_dir", "o", defaultOutputDir(), "Output directory")
	viper.BindPFlag("output_dir", rootCmd.PersistentFlags().Lookup("output_dir"))

	rootCmd.PersistentFlags().BoolP("serve_output", "s", DEFAULT_SERVE, "Start a server to serve built sites")
	viper.BindPFlag("serve_output", rootCmd.PersistentFlags().Lookup("serve_output"))

	rootCmd.PersistentFlags().IntP("serve_output_port", "T", DEFAULT_SERVE_PORT, "Port to serve built sites")
	viper.BindPFlag("serve_output_port", rootCmd.PersistentFlags().Lookup("serve_output_port"))

	// service
	rootCmd.PersistentFlags().String("service_name", DEFAULT_SERVICE_NAME, "Service name")
	viper.BindPFlag("service_name", rootCmd.PersistentFlags().Lookup("service_name"))

	rootCmd.PersistentFlags().String("service_logo", "", "Service logo image url")
	viper.BindPFlag("service_logo", rootCmd.PersistentFlags().Lookup("service_logo"))

	rootCmd.PersistentFlags().String("service_url", DEFAULT_SERVICE_URL, "Service URL")
	viper.BindPFlag("service_url", rootCmd.PersistentFlags().Lookup("service_url"))

	rootCmd.PersistentFlags().String("service_copyright_notice", DEFAULT_SERVICE_COPYRIGHT, "Service copyright notice")
	viper.BindPFlag("service_copyright_notice", rootCmd.PersistentFlags().Lookup("service_copyright_notice"))
}

func defaultOutputDir() string {
	return path.Join(helpers.WorkingDir(), DEFAULT_OUTPUT_DIR)
}

func checkAndOutputsFlags() {
	if viper.GetString("upload_dir") == "" {
		log.Fatalln("ERROR: The upload_dir setting is mandatory")
	}

	if viper.GetString("themes_dir") == "" {
		log.Fatalln("ERROR: The themes_dir setting is mandatory")
	}

	log.Printf("Upload dir: %s", viper.GetString("upload_dir"))
	log.Printf("Themes dir: %s", viper.GetString("themes_dir"))
	log.Printf("Output dir: %s", viper.GetString("output_dir"))
}

func setupConfig() {
	// setup environment variables
	// eg: KOWA_SMTP_AUTH_PASS
	viper.SetEnvPrefix("kowa")
	viper.AutomaticEnv()

	// setup config file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath("/etc/kowa/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.kowa/") // call multiple times to add many search paths
	viper.ReadInConfig()
}

// Add commands to root command
func addCommands() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(genDerivativesCmd)
	rootCmd.AddCommand(addUserCmd)
	rootCmd.AddCommand(addSiteCmd)
	rootCmd.AddCommand(fixImagesCmd)
	rootCmd.AddCommand(versionCmd)
}

//
// Main API
//

// Init commands configuration
func InitConf() {
	initKowaConf()
	initServerConf()
}

// Execute command
func Execute() {
	addCommands()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
