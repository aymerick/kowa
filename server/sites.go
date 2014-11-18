package server

import (
	"net/http"

	"github.com/aymerick/kowa/models"
)

// GET /sites
func listSites(w http.ResponseWriter, req *http.Request) {
	renderResp.JSON(w, http.StatusOK, renderMap{"sites": models.AllSites()})
}
