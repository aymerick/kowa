package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Run() {
	port := viper.GetString("port")

	r := gin.Default()
	setupRoutes(r)

	fmt.Println("Running on port:", port)
	r.Run(":" + port)
}

func setupRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/users", listUsers)
	r.GET("/sites", listSites)
}
