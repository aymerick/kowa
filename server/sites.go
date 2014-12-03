package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
)

// GET /api/sites/{site_id}
func (app *Application) handleGetSite(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSite\n")

	currentSite := context.Get(req, "currentSite").(*models.Site)
	if currentSite != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"site": currentSite})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /api/sites/{site_id}
func (app *Application) handleUpdateSite(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleUpdateSite\n")

	currentSite := context.Get(req, "currentSite").(*models.Site)
	if currentSite != nil {
		// @todo update site !
		panic("not implemented")

		app.render.JSON(rw, http.StatusOK, renderMap{"site": currentSite})
	} else {
		http.NotFound(rw, req)
	}
}
