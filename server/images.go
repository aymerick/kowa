package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
	"gopkg.in/mgo.v2/bson"
)

const (
	IMAGE_CT_PREFIX = "image/"
)

var AcceptedImageContentTypes []string

func init() {
	AcceptedImageContentTypes = []string{"image/jpeg", "image/png", "image/gif"}
}

// Check if given content type is an allowed image
func allowedImageContentType(ct string) bool {
	for _, allowedCT := range AcceptedImageContentTypes {
		if ct == allowedCT {
			return true
		}
	}

	return false
}

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

func (app *Application) handleGetImage(rw http.ResponseWriter, req *http.Request) {
	image := app.getCurrentImage(req)
	if image != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"image": image})
	} else {
		http.NotFound(rw, req)
	}
}

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

func (app *Application) handleUploadImage(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	site := app.getCurrentSite(req)
	if site == nil {
		panic("Site should be set")
	}

	reader, err := req.MultipartReader()
	if err != nil {
		log.Printf("Multipart error: %v", err.Error())
		http.Error(rw, "Failed to parse multipart data", http.StatusBadRequest)
		return
	}

	var fileName string
	var fileType string
	var fileInfo os.FileInfo

	for fileName == "" {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		fileName = part.FileName()
		if fileName == "" {
			continue
		}

		// Check content type
		fileType = part.Header.Get("Content-Type")

		if !allowedImageContentType(fileType) {
			log.Printf("Expected image content type but got: %v", fileType)
			http.Error(rw, "Unsupported image type", http.StatusBadRequest)
			return
		}

		fileType = strings.TrimPrefix(fileType, IMAGE_CT_PREFIX)

		// copy uploaded file
		log.Printf("Handling uploaded file: %s", fileName)

		dstPath := utils.AvailableFilePath(utils.AppUploadSiteFilePath(site.Id, fileName))

		dst, err := os.Create(dstPath)
		if err != nil {
			log.Printf("Can't create file: %s - %v", dstPath, err.Error())
			http.Error(rw, "Failed to create uploaded file", http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		if _, err := io.Copy(dst, part); err != nil {
			log.Printf("Can't save file: %s - %v", dstPath, err.Error())
			http.Error(rw, "Failed to save uploaded file", http.StatusInternalServerError)
			return
		}

		var errStat error
		fileInfo, errStat = os.Stat(dstPath)
		if os.IsNotExist(errStat) {
			http.Error(rw, "Failed to create uploaded file", http.StatusInternalServerError)
			return
		}
	}

	if fileName == "" {
		http.Error(rw, "Image not found in multipart", http.StatusBadRequest)
	} else {
		// create image model
		img := &models.Image{
			Id:     bson.NewObjectId(),
			SiteId: site.Id,
			Path:   fileInfo.Name(),
			Name:   fileName,
			Size:   fileInfo.Size(),
			Type:   fileType,
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
}
