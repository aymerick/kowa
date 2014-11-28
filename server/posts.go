package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
)

// GET /posts?site={site_id}
func (app *Application) handleGetPosts(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetPosts\n")

	site := context.Get(req, "currentSite").(*models.Site)
	if site != nil {
		skip, limit, err := paginationParams(req)
		if err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"posts": site.FindPosts(skip, limit)})
	} else {
		http.NotFound(rw, req)
	}
}
