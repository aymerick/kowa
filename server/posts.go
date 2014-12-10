package server

import (
	"log"
	"net/http"
)

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
