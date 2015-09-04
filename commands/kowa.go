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
	defaultMongodbURI    = "mongodb://localhost:27017/"
	defaultMongodbDbName = "kowa"
	defaultServe         = false
	defaultServePort     = 48910

	defaultServiceName      = "Kowa"
	defaultServiceURL       = "http://127.0.0.1" // @todo FIXME
	defaultServiceCopyright = "Copyright @ 2015 Kowa - All rights reserved"

	defaultOutputDir = "_sites"
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
	rootCmd.PersistentFlags().StringP("mongodb_uri", "d", defaultMongodbURI, "Uri to connect to mongoDB")
	viper.BindPFlag("mongodb_uri", rootCmd.PersistentFlags().Lookup("mongodb_uri"))

	rootCmd.PersistentFlags().StringP("mongodb_dbname", "n", defaultMongodbDbName, "MongoDB database name")
	viper.BindPFlag("mongodb_dbname", rootCmd.PersistentFlags().Lookup("mongodb_dbname"))

	// builder
	rootCmd.PersistentFlags().StringP("themes_dir", "t", "", "Themes directory (MANDATORY)")
	viper.BindPFlag("themes_dir", rootCmd.PersistentFlags().Lookup("themes_dir"))

	rootCmd.PersistentFlags().StringP("output_dir", "o", defaultOutputDirPath(), "Output directory")
	viper.BindPFlag("output_dir", rootCmd.PersistentFlags().Lookup("output_dir"))

	rootCmd.PersistentFlags().BoolP("serve_output", "s", defaultServe, "Start a server to serve built sites")
	viper.BindPFlag("serve_output", rootCmd.PersistentFlags().Lookup("serve_output"))

	rootCmd.PersistentFlags().IntP("serve_output_port", "T", defaultServePort, "Port to serve built sites")
	viper.BindPFlag("serve_output_port", rootCmd.PersistentFlags().Lookup("serve_output_port"))

	// service
	rootCmd.PersistentFlags().String("service_name", defaultServiceName, "Service name")
	viper.BindPFlag("service_name", rootCmd.PersistentFlags().Lookup("service_name"))

	rootCmd.PersistentFlags().String("service_logo", "", "Service logo image url")
	viper.BindPFlag("service_logo", rootCmd.PersistentFlags().Lookup("service_logo"))

	rootCmd.PersistentFlags().String("service_url", defaultServiceURL, "Service URL")
	viper.BindPFlag("service_url", rootCmd.PersistentFlags().Lookup("service_url"))

	rootCmd.PersistentFlags().String("service_copyright_notice", defaultServiceCopyright, "Service copyright notice")
	viper.BindPFlag("service_copyright_notice", rootCmd.PersistentFlags().Lookup("service_copyright_notice"))

	rootCmd.PersistentFlags().String("service_domains", "", "Service domains list")
	viper.BindPFlag("service_domains", rootCmd.PersistentFlags().Lookup("service_domains"))
}

func defaultOutputDirPath() string {
	return path.Join(helpers.WorkingDir(), defaultOutputDir)
}

func checkAndOutputsGlobalFlags() {
	if viper.GetString("upload_dir") == "" {
		log.Fatalln("ERROR: The upload_dir setting is mandatory")
	}

	if viper.GetString("themes_dir") == "" {
		log.Fatalln("ERROR: The themes_dir setting is mandatory")
	}

	if viper.GetString("service_url") == "" {
		log.Fatalln("ERROR: The service_url setting is mandatory")
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

// InitConf initializes commands configuration
func InitConf() {
	initKowaConf()
	initServerConf()
}

// Execute executes command
func Execute() {
	addCommands()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
