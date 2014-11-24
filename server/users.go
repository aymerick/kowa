package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// GET /api/users/{user_id}
func (app *Application) handleGetUser(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	if userId == "me" {
		// @todo Handle 'me' user
		userId = "5472f8ffc25c193af3000001"
	}

	// @todo Handle user not found
	user := app.dbSession.FindUser(userId)

	app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
}

// GET /api/users/{user_id}/sites
func (app *Application) handleGetUserSites(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	// @todo Handle NotFound
	user := app.dbSession.FindUser(userId)

	// @todo Check if user is currentUser

	app.render.JSON(rw, http.StatusOK, renderMap{"sites": user.FindSites()})
}
