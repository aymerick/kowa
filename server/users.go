package server

import (
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/mux"
)

type UsersController struct {
	*ApplicationController
}

// GET /api/users/{user_id}
func (this *UsersController) handleGetUser(rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	if userId == "me" {
		// @todo Handle 'me' user
		userId = "5472f8ffc25c193af3000001"
	}

	// @todo Handle user not found
	user := models.FindUser(userId)

	this.render.JSON(rw, http.StatusOK, renderMap{"user": user})

	return nil
}

// GET /api/users/{user_id}/sites
func (this *UsersController) handleGetUserSites(rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	// @todo Handle NotFound
	user := models.FindUser(userId)

	// @todo Check if user is currentUser

	this.render.JSON(rw, http.StatusOK, renderMap{"sites": user.FindSites()})

	return nil
}
