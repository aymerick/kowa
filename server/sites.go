package server

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// GET /api/users/{user_id}/sites/{site_id}
func (app *Application) handleGetSite(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSite\n")

	vars := mux.Vars(req)
	siteId := vars["site_id"]

	if site := app.dbSession.FindSite(bson.ObjectIdHex(siteId)); site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"site": site})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}/sites/{site_id}/posts
func (app *Application) handleGetSitePosts(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSitePosts\n")

	// @todo
}

// GET /api/users/{user_id}/sites/{site_id}/events
func (app *Application) handleGetSiteEvents(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSiteEvents\n")

	// @todo
}

// GET /api/users/{user_id}/sites/{site_id}/pages
func (app *Application) handleGetSitePages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSitePages\n")

	// @todo
}
