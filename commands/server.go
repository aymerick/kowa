package commands

import (
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/server"
)

const (
	DEFAULT_PORT = 35830

	DEFAULT_SMTP_FROM = "Kowa Server <kowa@localhost>"
	DEFAULT_SMTP_HOST = "127.0.0.1"
	DEFAULT_SMTP_PORT = 25
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Starts an HTTP server to handle API requests from the web app.`,
	Run:   runServer,
}

func initServerConf() {
	serverCmd.Flags().IntP("port", "p", DEFAULT_PORT, "Port to run Kowa server on")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))

	// SMTP

	serverCmd.Flags().String("smtp_from", DEFAULT_SMTP_FROM, "'from' email address to use when sending mails")
	viper.BindPFlag("smtp_from", serverCmd.Flags().Lookup("smtp_from"))

	serverCmd.Flags().String("smtp_host", DEFAULT_SMTP_HOST, "SMTP server host")
	viper.BindPFlag("smtp_host", serverCmd.Flags().Lookup("smtp_host"))

	serverCmd.Flags().Int("smtp_port", DEFAULT_SMTP_PORT, "SMTP server port")
	viper.BindPFlag("smtp_port", serverCmd.Flags().Lookup("smtp_port"))

	serverCmd.Flags().String("smtp_auth_user", "", "SMTP server username")
	viper.BindPFlag("smtp_auth_user", serverCmd.Flags().Lookup("smtp_auth_user"))

	serverCmd.Flags().String("smtp_auth_pass", "", "SMTP server password")
	viper.BindPFlag("smtp_auth_pass", serverCmd.Flags().Lookup("smtp_auth_pass"))
}

func runServer(cmd *cobra.Command, args []string) {
	checkAndOutputsFlags()

	app := server.NewApplication()
	app.Setup()

	go app.Run()

	// wait for interuption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	app.Stop()
}
