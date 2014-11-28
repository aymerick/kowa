package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// sugar helper
type renderMap map[string]interface{}

func Run() {
	// setup app
	app := NewApplication()

	// setup middlewares
	baseChain := alice.New(context.ClearHandler, app.loggingMiddleware, app.recoveryMiddleware, app.corsMiddleware())
	authChain := baseChain.Append(app.ensureAuthMiddleware)
	curUserChain := authChain.Append(app.ensureUserAccessMiddleware)

	// setup API routes
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// /api/oauth
	oauthRouter := apiRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.Methods("POST").Path("/token").Handler(baseChain.ThenFunc(app.handleOauthToken))
	oauthRouter.Methods("POST").Path("/revoke").Handler(baseChain.ThenFunc(app.handleOauthRevoke))

	// /api/me
	apiRouter.Methods("GET").Path("/me").Handler(authChain.ThenFunc(app.handleGetMe))

	// /api/users
	userRouter := apiRouter.PathPrefix("/users/{user_id}").Subrouter()
	userRouter.Methods("GET").Path("/sites/{site_id}/posts").Handler(curUserChain.ThenFunc(app.handleGetSitePosts))
	userRouter.Methods("GET").Path("/sites/{site_id}/events").Handler(curUserChain.ThenFunc(app.handleGetSiteEvents))
	userRouter.Methods("GET").Path("/sites/{site_id}/pages").Handler(curUserChain.ThenFunc(app.handleGetSitePages))
	userRouter.Methods("GET").Path("/sites/{site_id}/actions").Handler(curUserChain.ThenFunc(app.handleGetSiteActions))
	userRouter.Methods("GET").Path("/sites/{site_id}").Handler(curUserChain.ThenFunc(app.handleGetSite))
	userRouter.Methods("GET").Path("/sites").Handler(curUserChain.ThenFunc(app.handleGetUserSites))
	userRouter.Methods("GET").Handler(curUserChain.ThenFunc(app.handleGetUser))

	fmt.Println("Running on port:", app.port)
	http.ListenAndServe(":"+app.port, router)
}

func unauthorized(rw http.ResponseWriter) {
	http.Error(rw, "Not Authorized", http.StatusUnauthorized)
}
