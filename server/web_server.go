package server

import (
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/server/middlewares"
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

	// setup app
	app := NewApplication()

	commonMiddlewares := alice.New(middlewares.Logging, middlewares.Recovery, middlewares.Cors())

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// users
	userRouter := apiRouter.PathPrefix("/users/{user_id}").Subrouter()
	userRouter.Methods("GET").Path("/sites").Handler(commonMiddlewares.ThenFunc(app.handleGetUserSites))
	userRouter.Methods("GET").Handler(commonMiddlewares.ThenFunc(app.handleGetUser))

	// oauth
	oauthRouter := apiRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.Methods("POST").Path("/token").Handler(commonMiddlewares.ThenFunc(app.handleOauthToken))
	oauthRouter.Methods("POST").Path("/revoke").Handler(commonMiddlewares.ThenFunc(app.handleOauthRevoke))

	fmt.Println("Running on port:", port)
	http.ListenAndServe(":"+port, router)
}
