package server

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

// helper
type renderMap map[string]interface{}

// global renderer
var renderResp = render.New(render.Options{})

func Run() {
	port := viper.GetString("port")

	// setup server and middlewares
	n := negroni.Classic()
	setupMiddlewares(n)

	// setup router
	router := mux.NewRouter()
	setupRoutes(router)

	n.UseHandler(router)

	fmt.Println("Running on port:", port)
	n.Run(":" + port)
}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/users", listUsers).Methods("GET")
	router.HandleFunc("/sites", listSites).Methods("GET")
	router.HandleFunc("/token", oauthGetToken).Methods("POST")
}
