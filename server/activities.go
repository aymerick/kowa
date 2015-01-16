package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
)

// GET /activities?site={site_id}
// GET /sites/{site_id}/activities
func (app *Application) handleGetActivities(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetActivities\n")

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
