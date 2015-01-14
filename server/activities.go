package server

import (
	"log"
	"net/http"
)

// GET /activities?site={site_id}
// GET /sites/{site_id}/activities
func (app *Application) handleGetActivities(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetActivities\n")

	site := app.getCurrentSite(req)
	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"activities": site.FindActivities()})
	} else {
		http.NotFound(rw, req)
	}
}
