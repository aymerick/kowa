package commands

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgFile string

var RootCmd = &cobra.Command{
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

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $HOME/kowa/config.toml)")

	RootCmd.PersistentFlags().String("mongodb_uri", "mongodb://localhost:27017/", "Uri to connect to mongoDB")
	viper.BindPFlag("mongodb_uri", RootCmd.PersistentFlags().Lookup("mongodb_uri"))
}

func initConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath("/etc/kowa/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.kowa/") // call multiple times to add many search paths
	viper.ReadInConfig()
}

func addCommands() {
	RootCmd.AddCommand(serverCmd)
}

func Execute() {
	addCommands()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
