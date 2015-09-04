package commands

import (
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/server"
)

const (
	defaultPort = 35830

	defaultSMTPFrom = "Kowa Server <kowa@localhost>"
	defaultSMTPHost = "127.0.0.1"
	defaultSMTPPort = 25
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Starts an HTTP server to handle API requests from the web app.`,
	Run:   runServer,
}

func initServerConf() {
	serverCmd.Flags().IntP("port", "p", defaultPort, "Port to run Kowa server on")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))

	serverCmd.Flags().String("secret_key", "", "Secret key used to sign tokens")
	viper.BindPFlag("secret_key", serverCmd.Flags().Lookup("secret_key"))

	// Mail
	serverCmd.Flags().String("mail_tpl_dir", "", "Mail templates directory. If not provided, default templates are used.")
	viper.BindPFlag("mail_tpl_dir", serverCmd.Flags().Lookup("mail_tpl_dir"))

	serverCmd.Flags().String("smtp_from", defaultSMTPFrom, "'from' email address to use when sending mails")
	viper.BindPFlag("smtp_from", serverCmd.Flags().Lookup("smtp_from"))

	serverCmd.Flags().String("smtp_host", defaultSMTPHost, "SMTP server host")
	viper.BindPFlag("smtp_host", serverCmd.Flags().Lookup("smtp_host"))

	serverCmd.Flags().Int("smtp_port", defaultSMTPPort, "SMTP server port")
	viper.BindPFlag("smtp_port", serverCmd.Flags().Lookup("smtp_port"))

	serverCmd.Flags().String("smtp_auth_user", "", "SMTP server username")
	viper.BindPFlag("smtp_auth_user", serverCmd.Flags().Lookup("smtp_auth_user"))

	serverCmd.Flags().String("smtp_auth_pass", "", "SMTP server password")
	viper.BindPFlag("smtp_auth_pass", serverCmd.Flags().Lookup("smtp_auth_pass"))
}

func runServer(cmd *cobra.Command, args []string) {
	checkAndOutputsGlobalFlags()

	checkServerFlags()

	app := server.NewApplication()
	app.Setup()

	go app.Run()

	// wait for interuption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	app.Stop()
}

func checkServerFlags() {
	if viper.GetString("secret_key") == "" {
		log.Fatalln("ERROR: The secret_key setting is mandatory")
	}
}
