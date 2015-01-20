package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_MONGODB_URI    = "mongodb://localhost:27017/"
	DEFAULT_MONGODB_DBNAME = "kowa"
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
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(bootstrapCmd)
	rootCmd.AddCommand(resetCmd)
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
