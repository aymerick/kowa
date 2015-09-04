package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var acceptedFileContentTypes = []string{"application/pdf", "text/plain"}

// GET /files?site={site_id}
// GET /sites/{site_id}/files
func (app *Application) handleGetFiles(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		pagination := NewPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.FilesNb()

		app.render.JSON(rw, http.StatusOK, renderMap{"files": site.FindFiles(pagination.Skip, pagination.PerPage), "meta": pagination})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /files/{file_id}
func (app *Application) handleGetFile(rw http.ResponseWriter, req *http.Request) {
	file := app.getCurrentFile(req)
	if file != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"file": file})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /files/{file_id}
func (app *Application) handleDeleteFile(rw http.ResponseWriter, req *http.Request) {
	file := app.getCurrentFile(req)
	if file != nil {
		if err := file.Delete(); err != nil {
			http.Error(rw, "Failed to delete file", http.StatusInternalServerError)
		} else {
			// remove all references to file from site content
			site := app.getCurrentSite(req)

			if err := site.RemoveFileReferences(file); err != nil {
				log.Printf("Failed to remove file references: %v", err.Error())
				http.Error(rw, "Error while deleting file", http.StatusInternalServerError)
				return
			}

			// site content has changed
			app.onSiteChange(site)

			// returns deleted file
			app.render.JSON(rw, http.StatusOK, renderMap{"file": file})
		}
	} else {
		http.NotFound(rw, req)
	}
}

// POST /api/files/upload
func (app *Application) handleUploadFile(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	site := app.getCurrentSite(req)
	if site == nil {
		panic("Site should be set")
	}

	vars := mux.Vars(req)

	// get kind
	kind := vars["kind"]
	if kind == "" {
		panic("Kind should be set")
	}

	// check kind
	if !models.IsValidFileKind(kind) {
		log.Printf("Invalid file kind received: %v - Expected: %v", kind, models.FileKinds)
		http.Error(rw, "Invalid file kind", http.StatusBadRequest)
		return
	}

	// get uploaded file
	upload := handleUpload(rw, req, site, acceptedFileContentTypes)
	if upload == nil {
		// error is already handled by handleUpload() function
		return
	}

	// create file model
	f := &models.File{
		Id:     bson.NewObjectId(),
		SiteId: site.Id,
		Kind:   kind,
		Path:   upload.info.Name(),
		Name:   upload.name,
		Size:   upload.info.Size(),
		Type:   upload.ctype,
	}

	if err := currentDBSession.CreateFile(f); err != nil {
		log.Printf("Can't create record: %v - %v", f, err.Error())
		http.Error(rw, "Failed to create file record", http.StatusInternalServerError)
		return
	}

	switch kind {
	case models.FileMembership:
		site.SetMembership(f.Id)
	}

	// site content has changed
	app.onSiteChange(site)

	// returns uploaded file path
	app.render.JSON(rw, http.StatusCreated, renderMap{"file": f})
}
