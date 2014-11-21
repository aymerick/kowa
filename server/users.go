package server

import (
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/mux"
)

// GET /api/users/{user_id}
func handleGetUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	if userId == "me" {
		// @todo Handle 'me' user
		userId = "546b5830c25c19e01f000001"
	}

	// @todo Handle user not found
	user := models.FindUser(userId)

	renderResp.JSON(w, http.StatusOK, renderMap{"user": user})
}

// GET /api/users/{user_id}/sites
func handleGetUserSites(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["user_id"]

	// @todo Handle NotFound
	user := models.FindUser(userId)

	// @todo Check if user is currentUser

	renderResp.JSON(w, http.StatusOK, renderMap{"sites": user.FindSites()})
}
