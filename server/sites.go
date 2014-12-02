package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GET /api/sites/{site_id}
func (app *Application) handleGetSite(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetSite\n")

	vars := mux.Vars(req)
	siteId := vars["site_id"]

	if site := app.dbSession.FindSite(siteId); site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"site": site})
	} else {
		http.NotFound(rw, req)
	}
}
