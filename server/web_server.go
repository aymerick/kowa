package server

import (
	"fmt"

	"github.com/RangelReale/osin"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

// global oauth2 server
var oauthServer *osin.Server

// global renderer
var renderResp *render.Render

// sugar helper
type renderMap map[string]interface{}

func Run() {
	port := viper.GetString("port")

	// setup osin oauth2 server
	osinConfig := osin.NewServerConfig()
	osinConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.TOKEN}
	osinConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD, osin.REFRESH_TOKEN}

	oauthStorage := NewOAuthStorage()
	oauthServer = osin.NewServer(osinConfig, oauthStorage)

	// setup renderer
	renderResp = render.New(render.Options{})

	// setup routes
	n := negroni.Classic()

	setupCORSMiddleware(n)

	router := mux.NewRouter()
	router.HandleFunc("/users", handleUsers).Methods("GET")
	router.HandleFunc("/sites", handleSites).Methods("GET")

	oauthRouter := mux.NewRouter()
	oauthRouter.HandleFunc("/oauth/token", handleOauthToken).Methods("POST")

	router.Handle("/oauth/token", negroni.New(
		negroni.HandlerFunc(injectOauthSecretMiddleware),
		negroni.Wrap(oauthRouter),
	))

	n.UseHandler(router)

	fmt.Println("Running on port:", port)
	n.Run(":" + port)
}
