package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/aymerick/kowa/models"
)

type siteJson struct {
	Site models.Site `json:"site"`
}

type sitePageSettingsJson struct {
	SitePageSettings models.SitePageSettings `json:"sitePageSetting"`
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

// POST /sites/{site_id}/page-settings
// PUT /sites/{site_id}/page-settings/{setting_id}
func (app *Application) handleSetPageSettings(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	site := app.getCurrentSite(req)
	if site != nil {
		var respJson sitePageSettingsJson

		if err := json.NewDecoder(req.Body).Decode(&respJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		pageSettings := &respJson.SitePageSettings

		if vars["setting_id"] != "" {
			// this is an update
			existingPageSettings := site.PageSettings[pageSettings.Kind]
			if existingPageSettings == nil {
				http.NotFound(rw, req)
				return
			}

			existingIdStr := existingPageSettings.Id.Hex()

			badParamId := existingIdStr != vars["setting_id"]
			badContentId := (pageSettings.Id != "") && (existingIdStr != pageSettings.Id.Hex())

			if badParamId || badContentId {
				http.Error(rw, "Page settings id mismatch", http.StatusBadRequest)
				return
			} else if pageSettings.Id == "" {
				pageSettings.Id = existingPageSettings.Id
			}
		}

		err := site.SetPageSettings(pageSettings)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to add page settings", http.StatusInternalServerError)
			return
		}

		// site content has changed
		app.onSiteChange(site)

		app.render.JSON(rw, http.StatusOK, renderMap{"sitePageSetting": pageSettings})
	} else {
		http.NotFound(rw, req)
	}
}
