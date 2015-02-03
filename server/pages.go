package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
)

type pageJson struct {
	Page models.Page `json:"page"`
}

// GET /pages?site={site_id}
// GET /sites/{site_id}/pages
func (app *Application) handleGetPages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetPages\n")

	site := app.getCurrentSite(req)
	if site != nil {
		// fetch paginated records
		pagination := NewPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.PagesNb()

		pages := site.FindPages(pagination.Skip, pagination.PerPage)

		// fetch covers
		images := []*models.Image{}

		for _, page := range *pages {
			if image := page.FindCover(); image != nil {
				images = append(images, image)
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"pages": pages, "meta": pagination, "images": images})
	} else {
		http.NotFound(rw, req)
	}
}

// POST /pages
func (app *Application) handlePostPages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handlePostPages\n")

	currentDBSession := app.getCurrentDBSession(req)

	var reqJson pageJson

	if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	page := &reqJson.Page

	if page.SiteId == "" {
		http.Error(rw, "Missing site field in page record", http.StatusBadRequest)
		return
	}

	site := currentDBSession.FindSite(page.SiteId)
	if site == nil {
		http.Error(rw, "Site not found", http.StatusBadRequest)
		return
	}

	currentUser := app.getCurrentUser(req)
	if site.UserId != currentUser.Id {
		unauthorized(rw)
		return
	}

	if err := currentDBSession.CreatePage(page); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create page", http.StatusInternalServerError)
		return
	}

	// site content has changed
	app.onSiteChange(site)

	app.render.JSON(rw, http.StatusCreated, renderMap{"page": page})
}

// GET /pages/{page_id}
func (app *Application) handleGetPage(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetPage\n")

	page := app.getCurrentPage(req)
	if page != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"page": page})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /pages/{page_id}
func (app *Application) handleUpdatePage(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleUpdatePage\n")

	page := app.getCurrentPage(req)
	if page != nil {
		var reqJson pageJson

		if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		updated, err := page.Update(&reqJson.Page)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update page", http.StatusInternalServerError)
			return
		}

		if updated {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"page": page})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /pages/{page_id}
func (app *Application) handleDeletePage(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleDeletePage\n")

	page := app.getCurrentPage(req)
	if page != nil {
		if err := page.Delete(); err != nil {
			http.Error(rw, "Failed to delete page", http.StatusInternalServerError)
		} else {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)

			// returns deleted page
			app.render.JSON(rw, http.StatusOK, renderMap{"page": page})
		}
	} else {
		http.NotFound(rw, req)
	}
}
