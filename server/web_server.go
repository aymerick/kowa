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

	// setup API routes
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// /api/oauth
	oauthRouter := apiRouter.PathPrefix("/oauth").Subrouter()
	oauthRouter.Methods("POST").Path("/token").Handler(baseChain.ThenFunc(app.handleOauthToken))
	oauthRouter.Methods("POST").Path("/revoke").Handler(baseChain.ThenFunc(app.handleOauthRevoke))

	authChain := baseChain.Append(app.ensureAuthMiddleware)

	// /api/me
	apiRouter.Methods("GET").Path("/me").Handler(authChain.ThenFunc(app.handleGetMe))

	curUserChain := authChain.Append(app.ensureUserAccessMiddleware)

	// /api/users/{user_id}
	apiRouter.Methods("GET").Path("/users/{user_id}/sites").Handler(curUserChain.ThenFunc(app.handleGetUserSites))
	apiRouter.Methods("GET").Path("/users/{user_id}").Handler(curUserChain.ThenFunc(app.handleGetUser))

	curSiteOwnerChain := authChain.Append(app.ensureSiteMiddleware, app.ensureSiteOwnerAccessMiddleware)

	// /api/sites/{site_id}
	apiRouter.Methods("GET").Path("/sites/{site_id}/posts").Handler(curSiteOwnerChain.ThenFunc(app.handleGetSitePosts))
	apiRouter.Methods("GET").Path("/sites/{site_id}/events").Handler(curSiteOwnerChain.ThenFunc(app.handleGetSiteEvents))
	apiRouter.Methods("GET").Path("/sites/{site_id}/pages").Handler(curSiteOwnerChain.ThenFunc(app.handleGetSitePages))
	apiRouter.Methods("GET").Path("/sites/{site_id}/actions").Handler(curSiteOwnerChain.ThenFunc(app.handleGetSiteActions))
	apiRouter.Methods("GET").Path("/sites/{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleGetSite))

	fmt.Println("Running on port:", app.port)
	http.ListenAndServe(":"+app.port, router)
}

func unauthorized(rw http.ResponseWriter) {
	http.Error(rw, "Not Authorized", http.StatusUnauthorized)
}
