package server

import (
	"log"
	"net/http"
)

// GET /images?site={site_id}
// GET /sites/{site_id}/images
func (app *Application) handleGetImages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetImages\n")

	site := app.getCurrentSite(req)
	if site != nil {
		pagination := NewPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.ImagesNb()

		app.render.JSON(rw, http.StatusOK, renderMap{"images": site.FindImages(pagination.Skip, pagination.PerPage), "meta": pagination})
	} else {
		http.NotFound(rw, req)
	}
}

func (app *Application) handleGetImage(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetImage\n")

	image := app.getCurrentImage(req)
	if image != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"image": image})
	} else {
		http.NotFound(rw, req)
	}
}

func (app *Application) handleDeleteImage(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleDeleteImage\n")

	image := app.getCurrentImage(req)
	if image != nil {
		if err := image.Delete(); err != nil {
			http.Error(rw, "Failed to delete image", http.StatusInternalServerError)
		} else {
			app.render.JSON(rw, http.StatusOK, renderMap{"image": image})
		}
	} else {
		http.NotFound(rw, req)
	}
}
