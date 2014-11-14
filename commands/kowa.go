package commands

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_MONGODB_URI    = "mongodb://localhost:27017/"
	DEFAULT_MONGODB_DBNAME = "kowa"
)

var CfgFile string

var rootCmd = &cobra.Command{
	Use:   "kowa",
	Short: "Kowa generates a website for your association",
	Long:  `Koaw is a website generator that targets associations. It powers the asso.ninja web service.`,
	Run:   rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	go Server()

	// wait for interuption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

func initKowaConf() {
	cobra.OnInitialize(setupConfig)

	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $HOME/kowa/config.toml)")

	rootCmd.PersistentFlags().String("mongodb_uri", DEFAULT_MONGODB_URI, "Uri to connect to mongoDB")
	viper.BindPFlag("mongodb_uri", rootCmd.PersistentFlags().Lookup("mongodb_uri"))

	rootCmd.PersistentFlags().String("mongodb_dbname", DEFAULT_MONGODB_DBNAME, "MongoDB database name")
	viper.BindPFlag("mongodb_dbname", rootCmd.PersistentFlags().Lookup("mongodb_dbname"))
}

func setupConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath("/etc/kowa/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.kowa/") // call multiple times to add many search paths
	viper.ReadInConfig()
}

func addCommands() {
	rootCmd.AddCommand(serverCmd)
}

//
// Main API
//

func ResetConf() {
	rootCmd.ResetFlags()
	rootCmd.ResetCommands()

	serverCmd.ResetFlags()
	serverCmd.ResetCommands()

	viper.Reset()
}

func InitConf() {
	initKowaConf()
	initServerConf()
}

func Execute() {
	addCommands()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
