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
	log.Printf("[handler]: handleGetSite\n")

	currentSite := app.getCurrentSite(req)
	if currentSite != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"site": currentSite})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /api/sites/{site_id}
func (app *Application) handleUpdateSite(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleUpdateSite\n")

	currentSite := app.getCurrentSite(req)
	if currentSite != nil {
		var respJson siteJson

		if err := json.NewDecoder(req.Body).Decode(&respJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		if err := currentSite.Update(&respJson.Site); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update site", http.StatusInternalServerError)
			return
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"site": currentSite})
	} else {
		http.NotFound(rw, req)
	}
}
