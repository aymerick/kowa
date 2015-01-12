package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
)

type postJson struct {
	Post models.Post `json:"post"`
}

// GET /posts?site={site_id}
// GET /sites/{site_id}/posts
func (app *Application) handleGetPosts(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetPosts\n")

	site := app.getCurrentSite(req)
	if site != nil {
		// fetch paginated posts
		pagination := NewPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.PostsNb()

		posts := site.FindPosts(pagination.Skip, pagination.PerPage)

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
	log.Printf("[handler]: handlePostPosts\n")

	currentDBSession := app.getCurrentDBSession(req)

	var reqJson postJson

	if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	post := &reqJson.Post

	if post.SiteId == "" {
		http.Error(rw, "Missing site field in post record", http.StatusBadRequest)
		return
	}

	site := currentDBSession.FindSite(post.SiteId)
	if site == nil {
		http.Error(rw, "Site not found", http.StatusBadRequest)
		return
	}

	currentUser := app.getCurrentUser(req)
	if site.UserId != currentUser.Id {
		unauthorized(rw)
		return
	}

	if err := currentDBSession.CreatePost(post); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(rw, "Failed to create post", http.StatusInternalServerError)
		return
	}

	app.render.JSON(rw, http.StatusCreated, renderMap{"post": post})
}

// GET /posts/{post_id}
func (app *Application) handleGetPost(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetPost\n")

	post := app.getCurrentPost(req)
	if post != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"post": post})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /posts/{post_id}
func (app *Application) handleUpdatePost(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleUpdatePost\n")

	post := app.getCurrentPost(req)
	if post != nil {
		var reqJson postJson

		if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		if err := post.Update(&reqJson.Post); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update post", http.StatusInternalServerError)
			return
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"post": post})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /posts/{post_id}
func (app *Application) handleDeletePost(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleDeletePost\n")

	post := app.getCurrentPost(req)
	if post != nil {
		if err := post.Delete(); err != nil {
			http.Error(rw, "Failed to delete post", http.StatusInternalServerError)
		} else {
			// returns deleted post
			app.render.JSON(rw, http.StatusOK, renderMap{"post": post})
		}
	} else {
		http.NotFound(rw, req)
	}
}
