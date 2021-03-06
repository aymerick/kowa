package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
)

type postJSON struct {
	Post models.Post `json:"post"`
}

// GET /posts?site={site_id}
// GET /sites/{site_id}/posts
func (app *Application) handleGetPosts(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		// fetch paginated posts
		pagination := newPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.PostsNb()

		posts := site.FindPosts(pagination.Skip, pagination.PerPage, false)

		// fetch covers
		images := []*models.Image{}

		for _, post := range *posts {
			if image := post.FindCover(); image != nil {
				images = append(images, image)
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"posts": posts, "meta": pagination, "images": images})
	} else {
		http.NotFound(rw, req)
	}
}

// POST /posts
func (app *Application) handlePostPosts(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	var reqJSON postJSON

	if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// @todo [security] Check all fields !
	post := &reqJSON.Post

	if post.SiteID == "" {
		http.Error(rw, "Missing site field in post record", http.StatusBadRequest)
		return
	}

	site := currentDBSession.FindSite(post.SiteID)
	if site == nil {
		http.Error(rw, "Site not found", http.StatusBadRequest)
		return
	}

	currentUser := app.getCurrentUser(req)
	if site.UserID != currentUser.ID {
		unauthorized(rw)
		return
	}

	if err := currentDBSession.CreatePost(post); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// site content has changed
	app.onSiteChange(site)

	app.render.JSON(rw, http.StatusCreated, renderMap{"post": post})
}

// GET /posts/{post_id}
func (app *Application) handleGetPost(rw http.ResponseWriter, req *http.Request) {
	post := app.getCurrentPost(req)
	if post != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"post": post})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /posts/{post_id}
func (app *Application) handleUpdatePost(rw http.ResponseWriter, req *http.Request) {
	post := app.getCurrentPost(req)
	if post != nil {
		var reqJSON postJSON

		if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		// @todo [security] Check all fields !
		updated, err := post.Update(&reqJSON.Post)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update post", http.StatusInternalServerError)
			return
		}

		if updated {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"post": post})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /posts/{post_id}
func (app *Application) handleDeletePost(rw http.ResponseWriter, req *http.Request) {
	post := app.getCurrentPost(req)
	if post != nil {
		if err := post.Delete(); err != nil {
			http.Error(rw, "Failed to delete post", http.StatusInternalServerError)
		} else {
			site := app.getCurrentSite(req)

			// site content has changed
			app.onSiteChange(site)

			// returns deleted post
			app.render.JSON(rw, http.StatusOK, renderMap{"post": post})
		}
	} else {
		http.NotFound(rw, req)
	}
}
