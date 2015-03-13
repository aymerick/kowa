package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/aymerick/kowa/models"
)

type userJson struct {
	User models.User `json:"user"`
}

// GET /api/me
func (app *Application) handleGetMe(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	currentUser := app.getCurrentUser(req)
	userId := currentUser.Id

	if user := currentDBSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}
func (app *Application) handleGetUser(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	vars := mux.Vars(req)
	userId := vars["user_id"]

	if user := currentDBSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /api/users/{user_id}
func (app *Application) handleUpdateUser(rw http.ResponseWriter, req *http.Request) {
	user := app.getCurrentUser(req)
	if user != nil {
		var reqJson userJson

		if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		_, err := user.Update(&reqJson.User)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update user", http.StatusInternalServerError)
			return
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}/sites
func (app *Application) handleGetUserSites(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

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

	if user := currentDBSession.FindUser(userId); user != nil {
		var image *models.Image
		images := []*models.Image{}

		pageSettingsArray := []*models.SitePageSettings{}

		sites := user.FindSites()
		for _, site := range *sites {
			if image = site.FindLogo(); image != nil {
				images = append(images, image)
			}

			if image = site.FindCover(); image != nil {
				images = append(images, image)
			}

			for _, pageSettings := range site.PageSettings {
				pageSettingsArray = append(pageSettingsArray, pageSettings)

				if image = site.FindPageSettingsCover(pageSettings.Kind); image != nil {
					images = append(images, image)
				}
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"sites": sites, "images": images, "sitePageSettings": pageSettingsArray})
	} else {
		http.NotFound(rw, req)
	}
}
