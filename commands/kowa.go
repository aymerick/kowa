package commands

import (
	"fmt"
	"os"

	"github.com/aymerick/kowa/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_MONGODB_URI    = "mongodb://localhost:27017/"
	DEFAULT_MONGODB_DBNAME = "kowa"
	DEFAULT_SERVE          = false
	DEFAULT_SERVE_PORT     = 48910
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kowa",
	Short: "Kowa generadtes a website for your association",
	Long:  `Koaw is a website generator that targets associations. It powers the asso.ninja web service.`,
}

func initKowaConf() {
	cobra.OnInitialize(setupConfig)

	// config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kowa/config.toml)")

	// mongodb database
	rootCmd.PersistentFlags().StringP("mongodb_uri", "u", DEFAULT_MONGODB_URI, "Uri to connect to mongoDB")
	viper.BindPFlag("mongodb_uri", rootCmd.PersistentFlags().Lookup("mongodb_uri"))

	rootCmd.PersistentFlags().StringP("mongodb_dbname", "d", DEFAULT_MONGODB_DBNAME, "MongoDB database name")
	viper.BindPFlag("mongodb_dbname", rootCmd.PersistentFlags().Lookup("mongodb_dbname"))

	// builder
	rootCmd.PersistentFlags().StringP("working_dir", "w", utils.WorkingDir(), "Working directory")
	viper.BindPFlag("working_dir", rootCmd.PersistentFlags().Lookup("working_dir"))

	rootCmd.PersistentFlags().StringP("output_dir", "o", utils.DEFAULT_OUTPUT_DIR, "Output directory")
	viper.BindPFlag("output_dir", rootCmd.PersistentFlags().Lookup("output_dir"))

	rootCmd.PersistentFlags().BoolP("serve_output", "s", DEFAULT_SERVE, "Start a server to serve built sites")
	viper.BindPFlag("serve_output", rootCmd.PersistentFlags().Lookup("serve_output"))

	rootCmd.PersistentFlags().IntP("serve_output_port", "T", DEFAULT_SERVE_PORT, "Port to serve built sites")
	viper.BindPFlag("serve_output_port", rootCmd.PersistentFlags().Lookup("serve_output_port"))
}

func setupConfig() {
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
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(genDerivativesCmd)
}

//
// Main API
//

// Init commands configuration
func InitConf() {
	initKowaConf()
	initBuilderConf()
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
