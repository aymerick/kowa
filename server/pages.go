package server

import (
	"log"
	"net/http"
)

// GET /pages?site={site_id}
// GET /sites/{site_id}/pages
func (app *Application) handleGetPages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetPages\n")

	site := app.getCurrentSite(req)
	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"pages": site.FindPages()})
	} else {
		http.NotFound(rw, req)
	}
}
