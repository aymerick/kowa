package server

import (
	"log"
	"net/http"
)

// GET /actions?site={site_id}
// GET /sites/{site_id}/actions
func (app *Application) handleGetActions(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetActions\n")

	site := app.getCurrentSite(req)
	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"actions": site.FindActions()})
	} else {
		http.NotFound(rw, req)
	}
}
