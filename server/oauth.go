package server

import (
	"fmt"
	"net/http"
)

// POST /token
func oauthGetToken(w http.ResponseWriter, req *http.Request) {
	fmt.Println("rw: %v", w.Header())

	renderResp.JSON(w, http.StatusOK, renderMap{"token": "prout"})
}
