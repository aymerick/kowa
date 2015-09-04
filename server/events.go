package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
)

type eventJSON struct {
	Event models.Event `json:"event"`
}

// GET /events?site={site_id}
// GET /sites/{site_id}/events
func (app *Application) handleGetEvents(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		// fetch paginated events
		pagination := newPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.EventsNb()

		events := site.FindEvents(pagination.Skip, pagination.PerPage)

		// fetch covers
		images := []*models.Image{}

		for _, event := range *events {
			if image := event.FindCover(); image != nil {
				images = append(images, image)
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"events": events, "meta": pagination, "images": images})
	} else {
		http.NotFound(rw, req)
	}
}

// POST /events
func (app *Application) handlePostEvents(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	var reqJSON eventJSON

	if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// @todo [security] Check all fields !
	event := &reqJSON.Event

	if event.SiteID == "" {
		http.Error(rw, "Missing site field in event record", http.StatusBadRequest)
		return
	}

	site := currentDBSession.FindSite(event.SiteID)
	if site == nil {
		http.Error(rw, "Site not found", http.StatusBadRequest)
		return
	}

	currentUser := app.getCurrentUser(req)
	if site.UserID != currentUser.ID {
		unauthorized(rw)
		return
	}

	if err := currentDBSession.CreateEvent(event); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create event", http.StatusInternalServerError)
		return
	}

	// site content has changed
	app.onSiteChange(site)

	app.render.JSON(rw, http.StatusCreated, renderMap{"event": event})
}

// GET /events/{event_id}
func (app *Application) handleGetEvent(rw http.ResponseWriter, req *http.Request) {
	event := app.getCurrentEvent(req)
	if event != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"event": event})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /events/{event_id}
func (app *Application) handleUpdateEvent(rw http.ResponseWriter, req *http.Request) {
	event := app.getCurrentEvent(req)
	if event != nil {
		var reqJSON eventJSON

		if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		// @todo [security] Check all fields !
		updated, err := event.Update(&reqJSON.Event)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update event", http.StatusInternalServerError)
			return
		}

		if updated {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"event": event})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /events/{event_id}
func (app *Application) handleDeleteEvent(rw http.ResponseWriter, req *http.Request) {
	event := app.getCurrentEvent(req)
	if event != nil {
		if err := event.Delete(); err != nil {
			http.Error(rw, "Failed to delete event", http.StatusInternalServerError)
		} else {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)

			// returns deleted event
			app.render.JSON(rw, http.StatusOK, renderMap{"event": event})
		}
	} else {
		http.NotFound(rw, req)
	}
}
