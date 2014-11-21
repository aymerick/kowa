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
	osinConfig.AccessExpiration = 3600 // One hour
	osinConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.TOKEN}
	osinConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD, osin.REFRESH_TOKEN}
	osinConfig.ErrorStatusCode = 401

	oauthStorage := NewOAuthStorage()
	oauthServer = osin.NewServer(osinConfig, oauthStorage)

	// setup renderer
	renderResp = render.New(render.Options{})

	// setup routes
	n := negroni.Classic()

	setupCORSMiddleware(n)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/users/{user_id}", handleGetUser).Methods("GET")
	apiRouter.HandleFunc("/users/{user_id}/sites", handleGetUserSites).Methods("GET")

	oauthRouter := mux.NewRouter()
	// @todo Wtf ? It seems akwards
	oauthRouter.HandleFunc("/api/oauth/token", handleOauthToken).Methods("POST")
	oauthRouter.HandleFunc("/api/oauth/revoke", handleOauthRevoke).Methods("POST")

	// @todo Wtf ? It seems akwards
	apiRouter.Handle("/oauth/token", negroni.New(
		negroni.HandlerFunc(injectOauthSecretMiddleware),
		negroni.Wrap(oauthRouter),
	))

	// @todo Wtf ? It seems akwards
	apiRouter.Handle("/oauth/revoke", negroni.New(
		negroni.HandlerFunc(injectOauthSecretMiddleware),
		negroni.Wrap(oauthRouter),
	))

	n.UseHandler(apiRouter)

	fmt.Println("Running on port:", port)
	n.Run(":" + port)
}
