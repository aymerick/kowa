package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/models"
)

type siteJSON struct {
	Site models.Site `json:"site"`
}

type sitePageSettingsJSON struct {
	SitePageSettings models.SitePageSettings `json:"sitePageSetting"`
}

type siteThemeSettingsJSON struct {
	SiteThemeSettings models.SiteThemeSettings `json:"siteThemeSetting"`
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

	var reqJSON siteJSON

	if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// @todo [security] Check all fields !
	site := &reqJSON.Site

	T := i18n.MustTfunc(currentUser.Lang)

	// validate input
	errors := make(map[string]string)

	// check site name
	site.Name = strings.TrimSpace(site.Name)
	if site.Name == "" {
		errors["name"] = T("setup_missing_name")
	}

	// check site id
	site.ID = strings.TrimSpace(site.ID)
	if site.ID == "" {
		errors["id"] = T("setup_missing_id")
	}

	// check site id format
	if errors["id"] == "" {
		if site.ID != helpers.NormalizeToSiteID(site.ID) {
			errors["id"] = T("signup_id_invalid")
		}
	}

	// check site id length
	if errors["id"] == "" {
		if len(site.ID) < 3 {
			errors["id"] = T("signup_id_too_short")
		}
	}

	// check if site id is already taken
	if errors["id"] == "" {
		if exSite := currentDBSession.FindSite(site.ID); exSite != nil {
			errors["id"] = T("signup_id_not_available")
		}
	}

	// check site domain
	site.Domain = strings.TrimSpace(site.Domain)
	if site.Domain == "" {
		site.CustomURL = core.BaseUrl(site.ID)
	} else if !core.ValidDomain(site.Domain) {
		http.Error(rw, "Invalid domain provided", http.StatusBadRequest)
	}

	site.Tagline = strings.TrimSpace(site.Tagline)
	site.Description = strings.TrimSpace(site.Description)

	if len(errors) > 0 {
		app.render.JSON(rw, http.StatusBadRequest, renderMap{"errors": errors})
		return
	}

	site.UserID = currentUser.ID
	site.Lang = currentUser.Lang
	site.TZ = currentUser.TZ
	site.Theme = core.DefaultTheme
	site.NameInNavBar = true

	if err := currentDBSession.CreateSite(site); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create site", http.StatusInternalServerError)
		return
	}

	core.EnsureSiteUploadDir(site.ID)

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
		var respJSON siteJSON

		if err := json.NewDecoder(req.Body).Decode(&respJSON); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		prevBuildDir := site.BuildDir()

		// @todo [security] Check all fields !
		updated, err := site.Update(&respJSON.Site)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update site", http.StatusInternalServerError)
			return
		}

		if site.BuildDir() != prevBuildDir {
			app.deleteBuild(site, prevBuildDir)
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

// DELETE /sites/{site_id}
func (app *Application) handleDeleteSite(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		if err := site.Delete(); err != nil {
			http.Error(rw, "Failed to delete site", http.StatusInternalServerError)
			return
		}

		app.onSiteDeletion(site)

		// returns deleted site
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
		var respJSON sitePageSettingsJSON

		if err := json.NewDecoder(req.Body).Decode(&respJSON); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		// @todo [security] Check all fields !
		pageSettings := &respJSON.SitePageSettings

		if vars["setting_id"] != "" {
			// this is an update
			existingPageSettings := site.PageSettings[pageSettings.Kind]
			if existingPageSettings == nil {
				http.NotFound(rw, req)
				return
			}

			existingIDStr := existingPageSettings.ID.Hex()

			badParamID := existingIDStr != vars["setting_id"]
			badContentID := (pageSettings.ID != "") && (existingIDStr != pageSettings.ID.Hex())

			if badParamID || badContentID {
				http.Error(rw, "Page settings id mismatch", http.StatusBadRequest)
				return
			} else if pageSettings.ID == "" {
				pageSettings.ID = existingPageSettings.ID
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

// POST /sites/{site_id}/theme-settings
// PUT /sites/{site_id}/theme-settings/{setting_id}
func (app *Application) handleSetThemeSettings(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	site := app.getCurrentSite(req)
	if site != nil {
		var respJSON siteThemeSettingsJSON

		if err := json.NewDecoder(req.Body).Decode(&respJSON); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		// @todo [security] Check all fields !
		themeSettings := &respJSON.SiteThemeSettings

		if vars["setting_id"] != "" {
			// this is an update
			existingThemeSettings := site.ThemeSettings[themeSettings.Theme]
			if existingThemeSettings == nil {
				http.NotFound(rw, req)
				return
			}

			existingIDStr := existingThemeSettings.ID.Hex()

			badParamID := existingIDStr != vars["setting_id"]
			badContentID := (themeSettings.ID != "") && (existingIDStr != themeSettings.ID.Hex())

			if badParamID || badContentID {
				http.Error(rw, "Theme settings id mismatch", http.StatusBadRequest)
				return
			} else if themeSettings.ID == "" {
				themeSettings.ID = existingThemeSettings.ID
			}
		}

		err := site.SetThemeSettings(themeSettings)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to add theme settings", http.StatusInternalServerError)
			return
		}

		// site content has changed
		app.onSiteChange(site)

		app.render.JSON(rw, http.StatusOK, renderMap{"siteThemeSetting": themeSettings})
	} else {
		http.NotFound(rw, req)
	}
}
