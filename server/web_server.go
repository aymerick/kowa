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

	// setup controllers
	app := &ApplicationController{render: render.New(render.Options{})}
	users := &UsersController{app}
	oauth := &OauthController{app}

	// setup routes
	n := negroni.Classic()

	setupCORSMiddleware(n)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// UsersController
	userRouter := apiRouter.PathPrefix("/users/{user_id}").Subrouter()
	userRouter.Methods("GET").Path("/sites").Handler(users.Action(users.handleGetUserSites))
	userRouter.Methods("GET").Handler(users.Action(users.handleGetUser))

	// OauthController
	oauthRouter := apiRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.Methods("POST").Path("/token").Handler(oauth.Action(oauth.handleOauthToken))
	oauthRouter.Methods("POST").Path("/revoke").Handler(oauth.Action(oauth.handleOauthRevoke))

	n.UseHandler(router)

	fmt.Println("Running on port:", port)
	n.Run(":" + port)
}
