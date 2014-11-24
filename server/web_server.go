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

	// /api/users/{user_id}
	userRouter := apiRouter.PathPrefix("/users/{user_id}").Subrouter()
	userRouter.Methods("GET").Path("/sites").HandlerFunc(handleGetUserSites)
	userRouter.Methods("GET").HandlerFunc(handleGetUser)

	// /api/oauth
	oauthRouter := apiRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.Methods("POST").Path("/token").HandlerFunc(handleOauthToken)
	oauthRouter.Methods("POST").Path("/revoke").HandlerFunc(handleOauthRevoke)

	n.UseHandler(router)

	fmt.Println("Running on port:", port)
	n.Run(":" + port)
}
