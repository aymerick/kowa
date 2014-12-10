package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GET /api/me
func (app *Application) handleGetMe(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetMe\n")

	currentUser := app.getCurrentUser(req)
	userId := currentUser.Id

	if user := app.dbSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}
func (app *Application) handleGetUser(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetUser\n")

	vars := mux.Vars(req)
	userId := vars["user_id"]

	if user := app.dbSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}/sites
func (app *Application) handleGetUserSites(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetUserSites\n")

	vars := mux.Vars(req)
	userId := vars["user_id"]

	// check current user
	currentUser := app.getCurrentUser(req)
	if currentUser == nil {
		unauthorized(rw)
		return
	}

	if currentUser.Id != userId {
		unauthorized(rw)
		return
	}

	if user := app.dbSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"sites": user.FindSites()})
	} else {
		http.NotFound(rw, req)
	}
}
