package server

import (
	"log"
	"net/http"

	"github.com/aymerick/kowa/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	IMAGE_CT_PREFIX = "image/"
)

var acceptedImageContentTypes = []string{"image/jpeg", "image/png", "image/gif"}

// GET /images?site={site_id}
// GET /sites/{site_id}/images
func (app *Application) handleGetImages(rw http.ResponseWriter, req *http.Request) {
	site := app.getCurrentSite(req)
	if site != nil {
		pagination := NewPagination()
		if err := pagination.fillFromRequest(req); err != nil {
			http.Error(rw, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		pagination.Total = site.ImagesNb()

		app.render.JSON(rw, http.StatusOK, renderMap{"images": site.FindImages(pagination.Skip, pagination.PerPage), "meta": pagination})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /images/{image_id}
func (app *Application) handleGetImage(rw http.ResponseWriter, req *http.Request) {
	image := app.getCurrentImage(req)
	if image != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"image": image})
	} else {
		http.NotFound(rw, req)
	}
}

// DELETE /images/{image_id}
func (app *Application) handleDeleteImage(rw http.ResponseWriter, req *http.Request) {
	image := app.getCurrentImage(req)
	if image != nil {
		if err := image.Delete(); err != nil {
			http.Error(rw, "Failed to delete image", http.StatusInternalServerError)
		} else {
			// remove all references to image from site content
			site := app.getCurrentSite(req)

			if err := site.RemoveImageReferences(image); err != nil {
				log.Printf("Failed to remove image references: %v", err.Error())
				http.Error(rw, "Error while deleting image", http.StatusInternalServerError)
				return
			}

			// site content has changed
			app.onSiteChange(site)

			// returns deleted image
			app.render.JSON(rw, http.StatusOK, renderMap{"image": image})
		}
	} else {
		http.NotFound(rw, req)
	}
}

// POST /images/upload
func (app *Application) handleUploadImage(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	site := app.getCurrentSite(req)
	if site == nil {
		panic("Site should be set")
	}

	// get uploaded file
	upload := handleUpload(rw, req, site, acceptedImageContentTypes)
	if upload == nil {
		// error is already handled by handleUpload() function
		return
	}

	// create image model
	img := &models.Image{
		Id:     bson.NewObjectId(),
		SiteId: site.Id,
		Path:   upload.info.Name(),
		Name:   upload.name,
		Size:   upload.info.Size(),
		Type:   upload.ctype,
	}

	if err := currentDBSession.CreateImage(img); err != nil {
		log.Printf("Can't create record: %v - %v", img, err.Error())
		http.Error(rw, "Failed to create image record", http.StatusInternalServerError)
		return
	}

	// @todo Async that the day it becomes problematic
	if err := img.GenerateDerivatives(true); err != nil {
		log.Printf("Failed to generate image derivatives: %s - %v", img.Path, err.Error())
	}

	// returns uploaded file path
	app.render.JSON(rw, http.StatusCreated, renderMap{"image": img})
}
