package server

import (
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/mux"
)

// GET /users/{user_id}
func handleGetUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// @todo Handle NotFound
	user := models.FindUser(vars["user_id"])

	renderResp.JSON(w, http.StatusOK, renderMap{"user": user})
}
