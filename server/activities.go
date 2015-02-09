package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
)

type activityJson struct {
	Activity models.Activity `json:"activity"`
}

// GET /activities?site={site_id}
// GET /sites/{site_id}/activities
func (app *Application) handleGetActivities(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		// fetch paginated records
		pagination := NewPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.ActivitiesNb()

		activities := site.FindActivities(pagination.Skip, pagination.PerPage)

		// fetch covers
		images := []*models.Image{}

		for _, activity := range *activities {
			if image := activity.FindCover(); image != nil {
				images = append(images, image)
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"activities": activities, "meta": pagination, "images": images})
	} else {
		http.NotFound(rw, req)
	}
}

// POST /activities
func (app *Application) handlePostActivities(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	var reqJson activityJson

	if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	activity := &reqJson.Activity

	if activity.SiteId == "" {
		http.Error(rw, "Missing site field in activity record", http.StatusBadRequest)
		return
	}

	site := currentDBSession.FindSite(activity.SiteId)
	if site == nil {
		http.Error(rw, "Site not found", http.StatusBadRequest)
		return
	}

	currentUser := app.getCurrentUser(req)
	if site.UserId != currentUser.Id {
		unauthorized(rw)
		return
	}

	if err := currentDBSession.CreateActivity(activity); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create activity", http.StatusInternalServerError)
		return
	}

	// site content has changed
	app.onSiteChange(site)

	app.render.JSON(rw, http.StatusCreated, renderMap{"activity": activity})
}

// GET /activities/{activity_id}
func (app *Application) handleGetActivity(rw http.ResponseWriter, req *http.Request) {
	activity := app.getCurrentActivity(req)
	if activity != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"activity": activity})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /activities/{activity_id}
func (app *Application) handleUpdateActivity(rw http.ResponseWriter, req *http.Request) {
	activity := app.getCurrentActivity(req)
	if activity != nil {
		var reqJson activityJson

		if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		updated, err := activity.Update(&reqJson.Activity)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update activity", http.StatusInternalServerError)
			return
		}

		if updated {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"activity": activity})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /activities/{activity_id}
func (app *Application) handleDeleteActivity(rw http.ResponseWriter, req *http.Request) {
	activity := app.getCurrentActivity(req)
	if activity != nil {
		if err := activity.Delete(); err != nil {
			http.Error(rw, "Failed to delete activity", http.StatusInternalServerError)
		} else {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)

			// returns deleted activity
			app.render.JSON(rw, http.StatusOK, renderMap{"activity": activity})
		}
	} else {
		http.NotFound(rw, req)
	}
}
