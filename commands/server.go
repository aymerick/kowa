package commands

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
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
	Run:   serverRun,
}

func initServerConf() {
	serverCmd.Flags().Int("port", DEFAULT_PORT, "Port to run Kowa server on")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
}

func serverRun(cmd *cobra.Command, args []string) {
	go Server()

	// wait for interuption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

func Server() {
	port := viper.GetString("port")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	fmt.Println("Running on port:", port)
	r.Run(":" + port)
}
