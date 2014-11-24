package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// sugar helper
type renderMap map[string]interface{}

func Run() {
	// setup app
	app := NewApplication()

	// setup middlewares
	baseChain := alice.New(app.loggingMiddleware, app.recoveryMiddleware, app.corsMiddleware())
	authChain := baseChain.Append(app.setCurrentUserMiddleware)

	// setup routes
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// users
	userRouter := apiRouter.PathPrefix("/users/{user_id}").Subrouter()
	userRouter.Methods("GET").Path("/sites").Handler(authChain.ThenFunc(app.handleGetUserSites))
	userRouter.Methods("GET").Handler(authChain.ThenFunc(app.handleGetUser))

	// oauth
	oauthRouter := apiRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.Methods("POST").Path("/token").Handler(baseChain.ThenFunc(app.handleOauthToken))
	oauthRouter.Methods("POST").Path("/revoke").Handler(baseChain.ThenFunc(app.handleOauthRevoke))

	fmt.Println("Running on port:", app.port)
	http.ListenAndServe(":"+app.port, router)
}
