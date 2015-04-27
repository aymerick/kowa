package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/models"
)

type siteJson struct {
	Site models.Site `json:"site"`
}

type sitePageSettingsJson struct {
	SitePageSettings models.SitePageSettings `json:"sitePageSetting"`
}

// POST /api/sites
func (app *Application) handlePostSite(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	currentUser := app.getCurrentUser(req)

	// @todo Check if user is allowed to create a new site
	// if !currentUser.canCreateSite() {
	// 	unauthorized(rw)
	// 	return
	// }

	var reqJson siteJson

	if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// @todo [security] Check all fields !
	site := &reqJson.Site

	T := i18n.MustTfunc(currentUser.Lang)

	// validate input
	errors := make(map[string]string)

	// check site name
	site.Name = strings.TrimSpace(site.Name)
	if site.Name == "" {
		errors["name"] = T("setup_missing_name")
	}

	// check site id
	site.Id = strings.TrimSpace(site.Id)
	if site.Id == "" {
		errors["id"] = T("setup_missing_id")
	}

	// check site id format
	if site.Id != helpers.NormalizeToPathPart(site.Id) {
		errors["id"] = T("signup_id_invalid")
	}

	// check if site id is already taken
	if exSite := currentDBSession.FindSite(site.Id); exSite != nil {
		errors["id"] = T("setup_id_not_available")
	}

	site.Tagline = strings.TrimSpace(site.Tagline)
	site.Description = strings.TrimSpace(site.Description)

	if len(errors) > 0 {
		app.render.JSON(rw, http.StatusBadRequest, renderMap{"errors": errors})
		return
	}

	site.UserId = currentUser.Id
	site.Lang = currentUser.Lang
	site.Theme = core.DEFAULT_THEME
	site.NameInNavBar = true

	// @todo FIXME !
	site.BaseUrl = fmt.Sprintf("%s:%d/%s", core.DEFAULT_BASEURL, viper.GetInt("serve_output_port"), site.Id)

	if err := currentDBSession.CreateSite(site); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create site", http.StatusInternalServerError)
		return
	}

	core.EnsureSiteUploadDir(site.Id)

	// site content has changed
	app.onSiteChange(site)

	app.render.JSON(rw, http.StatusCreated, renderMap{"site": site})
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

		// @todo [security] Check all fields !
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

		// @todo [security] Check all fields !
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
