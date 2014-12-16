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
		pagination := NewPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.PostsNb()

		app.render.JSON(rw, http.StatusOK, renderMap{"posts": site.FindPosts(pagination.Skip, pagination.PerPage), "meta": pagination})
	} else {
		http.NotFound(rw, req)
	}
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
		var respJson postJson

		if err := json.NewDecoder(req.Body).Decode(&respJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		if err := post.Update(&respJson.Post); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update post", http.StatusInternalServerError)
			return
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"post": post})
	} else {
		http.NotFound(rw, req)
	}
}
