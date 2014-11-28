package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
)

// GET /pages?site={site_id}
func (app *Application) handleGetPages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetPages\n")

	site := context.Get(req, "currentSite").(*models.Site)

	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"pages": site.FindPages()})
	} else {
		http.NotFound(rw, req)
	}
}
