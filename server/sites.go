package server

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// GET /api/sites/{site_id}
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

// GET /api/sites/{site_id}/posts
func (app *Application) handleGetSitePosts(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSitePosts\n")

	site := context.Get(req, "currentSite").(*models.Site)

	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"posts": site.FindPosts()})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/sites/{site_id}/events
func (app *Application) handleGetSiteEvents(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSiteEvents\n")

	site := context.Get(req, "currentSite").(*models.Site)

	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"events": site.FindEvents()})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/sites/{site_id}/pages
func (app *Application) handleGetSitePages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSitePages\n")

	site := context.Get(req, "currentSite").(*models.Site)

	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"pages": site.FindPages()})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/sites/{site_id}/actions
func (app *Application) handleGetSiteActions(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSiteActions\n")

	site := context.Get(req, "currentSite").(*models.Site)

	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"actions": site.FindActions()})
	} else {
		http.NotFound(rw, req)
	}
}
