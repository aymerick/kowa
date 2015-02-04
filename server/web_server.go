package server

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// sugar helper
type renderMap map[string]interface{}

func (app *Application) newWebRouter() *mux.Router {
	// setup middlewares
	baseChain := alice.New(context.ClearHandler, app.dbSessionMiddleware, app.loggingMiddleware, app.recoveryMiddleware, app.corsMiddleware())

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

	// midllewares
	curSiteOwnerChain := authChain.Append(app.ensureSiteMiddleware, app.ensureSiteOwnerAccessMiddleware)
	curPostOwnerChain := authChain.Append(app.ensurePostMiddleware, app.ensureSiteMiddleware, app.ensureSiteOwnerAccessMiddleware)
	curPageOwnerChain := authChain.Append(app.ensurePageMiddleware, app.ensureSiteMiddleware, app.ensureSiteOwnerAccessMiddleware)
	curActivityOwnerChain := authChain.Append(app.ensureActivityMiddleware, app.ensureSiteMiddleware, app.ensureSiteOwnerAccessMiddleware)
	curImageOwnerChain := authChain.Append(app.ensureImageMiddleware, app.ensureSiteMiddleware, app.ensureSiteOwnerAccessMiddleware)

	// /api/sites/{site_id}
	apiRouter.Methods("GET").Path("/sites/{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleGetSite))
	apiRouter.Methods("PUT").Path("/sites/{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleUpdateSite))
	apiRouter.Methods("GET").Path("/sites/{site_id}/posts").Handler(curSiteOwnerChain.ThenFunc(app.handleGetPosts))
	apiRouter.Methods("GET").Path("/sites/{site_id}/events").Handler(curSiteOwnerChain.ThenFunc(app.handleGetEvents))
	apiRouter.Methods("GET").Path("/sites/{site_id}/pages").Handler(curSiteOwnerChain.ThenFunc(app.handleGetPages))
	apiRouter.Methods("GET").Path("/sites/{site_id}/activities").Handler(curSiteOwnerChain.ThenFunc(app.handleGetActivities))
	apiRouter.Methods("GET").Path("/sites/{site_id}/images").Handler(curSiteOwnerChain.ThenFunc(app.handleGetImages))

	// /api/posts?site={site_id}
	apiRouter.Methods("GET").Path("/posts").Queries("site", "{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleGetPosts))
	apiRouter.Methods("POST").Path("/posts").Handler(authChain.ThenFunc(app.handlePostPosts))
	apiRouter.Methods("GET").Path("/posts/{post_id}").Handler(curPostOwnerChain.ThenFunc(app.handleGetPost))
	apiRouter.Methods("PUT").Path("/posts/{post_id}").Handler(curPostOwnerChain.ThenFunc(app.handleUpdatePost))
	apiRouter.Methods("DELETE").Path("/posts/{post_id}").Handler(curPostOwnerChain.ThenFunc(app.handleDeletePost))

	// /api/events?site={site_id}
	apiRouter.Methods("GET").Path("/events").Queries("site", "{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleGetEvents))

	// /api/pages?site={site_id}
	apiRouter.Methods("GET").Path("/pages").Queries("site", "{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleGetPages))
	apiRouter.Methods("POST").Path("/pages").Handler(authChain.ThenFunc(app.handlePostPages))
	apiRouter.Methods("GET").Path("/pages/{page_id}").Handler(curPageOwnerChain.ThenFunc(app.handleGetPage))
	apiRouter.Methods("PUT").Path("/pages/{page_id}").Handler(curPageOwnerChain.ThenFunc(app.handleUpdatePage))
	apiRouter.Methods("DELETE").Path("/pages/{page_id}").Handler(curPageOwnerChain.ThenFunc(app.handleDeletePage))

	// /api/activities?site={site_id}
	apiRouter.Methods("GET").Path("/activities").Queries("site", "{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleGetActivities))
	apiRouter.Methods("POST").Path("/activities").Handler(authChain.ThenFunc(app.handlePostActivities))
	apiRouter.Methods("GET").Path("/activities/{activity_id}").Handler(curActivityOwnerChain.ThenFunc(app.handleGetActivity))
	apiRouter.Methods("PUT").Path("/activities/{activity_id}").Handler(curActivityOwnerChain.ThenFunc(app.handleUpdateActivity))
	apiRouter.Methods("DELETE").Path("/activities/{activity_id}").Handler(curActivityOwnerChain.ThenFunc(app.handleDeleteActivity))

	// /api/images?site={site_id}
	apiRouter.Methods("GET").Path("/images").Queries("site", "{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleGetImages))
	apiRouter.Methods("GET").Path("/images/{image_id}").Handler(curImageOwnerChain.ThenFunc(app.handleGetImage))
	apiRouter.Methods("DELETE").Path("/images/{image_id}").Handler(curImageOwnerChain.ThenFunc(app.handleDeleteImage))

	// upload image
	apiRouter.Methods("POST").Path("/images/upload").Queries("site", "{site_id}").Handler(curSiteOwnerChain.ThenFunc(app.handleUploadImage))

	return router
}

func unauthorized(rw http.ResponseWriter) {
	http.Error(rw, "Not Authorized", http.StatusUnauthorized)
}
