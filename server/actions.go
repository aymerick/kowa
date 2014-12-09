package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
)

// GET /actions?site={site_id}
// GET /sites/{site_id}/actions
func (app *Application) handleGetActions(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetActions\n")

	site := context.Get(req, "currentSite").(*models.Site)

	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"actions": site.FindActions()})
	} else {
		http.NotFound(rw, req)
	}
}
