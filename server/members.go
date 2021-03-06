package server

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/aymerick/kowa/models"
)

type memberJSON struct {
	Member models.Member `json:"member"`
}

// GET /members?site={site_id}
// GET /sites/{site_id}/members
func (app *Application) handleGetMembers(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		// fetch paginated records
		pagination := newPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.MembersNb()

		members := site.FindMembers(pagination.Skip, pagination.PerPage)

		// fetch photos
		images := []*models.Image{}

		for _, member := range *members {
			if image := member.FindPhoto(); image != nil {
				images = append(images, image)
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"members": members, "meta": pagination, "images": images})
	} else {
		http.NotFound(rw, req)
	}
}

// POST /members
func (app *Application) handlePostMembers(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	var reqJSON memberJSON

	if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// @todo [security] Check all fields !
	member := &reqJSON.Member

	if member.SiteID == "" {
		http.Error(rw, "Missing site field in member record", http.StatusBadRequest)
		return
	}

	site := currentDBSession.FindSite(member.SiteID)
	if site == nil {
		http.Error(rw, "Site not found", http.StatusBadRequest)
		return
	}

	currentUser := app.getCurrentUser(req)
	if site.UserID != currentUser.ID {
		unauthorized(rw)
		return
	}

	if err := currentDBSession.CreateMember(member); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create member", http.StatusInternalServerError)
		return
	}

	// site content has changed
	app.onSiteChange(site)

	app.render.JSON(rw, http.StatusCreated, renderMap{"member": member})
}

// GET /members/{member_id}
func (app *Application) handleGetMember(rw http.ResponseWriter, req *http.Request) {
	member := app.getCurrentMember(req)
	if member != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"member": member})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /members/{member_id}
func (app *Application) handleUpdateMember(rw http.ResponseWriter, req *http.Request) {
	member := app.getCurrentMember(req)
	if member != nil {
		var reqJSON memberJSON

		if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		// @todo [security] Check all fields !
		updated, err := member.Update(&reqJSON.Member)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update member", http.StatusInternalServerError)
			return
		}

		if updated {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"member": member})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /members/{member_id}
func (app *Application) handleDeleteMember(rw http.ResponseWriter, req *http.Request) {
	member := app.getCurrentMember(req)
	if member != nil {
		if err := member.Delete(); err != nil {
			http.Error(rw, "Failed to delete member", http.StatusInternalServerError)
		} else {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)

			// returns deleted member
			app.render.JSON(rw, http.StatusOK, renderMap{"member": member})
		}
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /members/order
func (app *Application) handlePutMembersOrder(rw http.ResponseWriter, req *http.Request) {
	var ids []bson.ObjectId

	if err := json.NewDecoder(req.Body).Decode(&ids); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	site := app.getCurrentSite(req)

	order := 1
	for _, id := range ids {
		site.UpdateMemberOrder(id, order)
		order++
	}

	// site content has changed
	app.onSiteChange(site)

	app.render.JSON(rw, http.StatusOK, renderMap{})
}
