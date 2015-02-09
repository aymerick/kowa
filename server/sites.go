package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
)

type siteJson struct {
	Site models.Site `json:"site"`
}

// GET /api/sites/{site_id}
func (app *Application) handleGetSite(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"site": site})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /api/sites/{site_id}
func (app *Application) handleUpdateSite(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		var respJson siteJson

		if err := json.NewDecoder(req.Body).Decode(&respJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		updated, err := site.Update(&respJson.Site)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update site", http.StatusInternalServerError)
			return
		}

		if updated {
			// site content has changed
			app.onSiteChange(site)
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"site": site})
	} else {
		http.NotFound(rw, req)
	}
}
