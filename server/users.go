package server

import (
	"net/http"

	"github.com/aymerick/kowa/models"
)

// GET /users
func handleUsers(w http.ResponseWriter, req *http.Request) {
	renderResp.JSON(w, http.StatusOK, renderMap{"users": models.AllUsers()})
}
