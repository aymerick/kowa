package server

import (
	"log"
	"net/http"
)

// GET /events?site={site_id}
// GET /sites/{site_id}/events
func (app *Application) handleGetEvents(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetEvents\n")

	site := app.getCurrentSite(req)
	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"events": site.FindEvents()})
	} else {
		http.NotFound(rw, req)
	}
}
