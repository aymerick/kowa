package server

import (
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// GET /api/users/{user_id}
func (app *Application) handleGetUser(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	if userId == "me" {
		currentUser := context.Get(req, "currentUser").(*models.User)
		userId = currentUser.Id
	}

	if user := app.dbSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}/sites
func (app *Application) handleGetUserSites(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	// check current user
	currentUser := context.Get(req, "currentUser").(*models.User)
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
