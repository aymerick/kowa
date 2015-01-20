package commands

import (
	"os"
	"os/signal"

	"github.com/aymerick/kowa/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DEFAULT_PORT = 35830
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server for admin API",
	Long:  `Kowa will starts an HTTP server to handle API requests from the web client.`,
	Run:   runServer,
}

func initServerConf() {
	serverCmd.Flags().IntP("port", "p", DEFAULT_PORT, "Port to run Kowa server on")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
}

func runServer(cmd *cobra.Command, args []string) {
	go server.Run()

	// wait for interuption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
