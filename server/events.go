package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
)

// GET /events?site={site_id}
func (app *Application) handleGetEvents(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetEvents\n")

	site := context.Get(req, "currentSite").(*models.Site)

	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"events": site.FindEvents()})
	} else {
		http.NotFound(rw, req)
	}
}
