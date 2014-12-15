package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"
	"gopkg.in/mgo.v2/bson"
)

const (
	PUBLIC_UPLOAD_PATH = "/upload"
)

// @todo FIXME
var appUploadDir string

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	appUploadDir = path.Join(currentDir, "/client/public", PUBLIC_UPLOAD_PATH)
}

// GET /images?site={site_id}
// GET /sites/{site_id}/images
func (app *Application) handleGetImages(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleGetImages\n")

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
	log.Printf("[handler]: handleGetImage\n")

	image := app.getCurrentImage(req)
	if image != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"image": image})
	} else {
		http.NotFound(rw, req)
	}
}

func (app *Application) handleDeleteImage(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleDeleteImage\n")

	image := app.getCurrentImage(req)
	if image != nil {
		if err := image.Delete(); err != nil {
			http.Error(rw, "Failed to delete image", http.StatusInternalServerError)
		} else {
			// update site
			site := app.getCurrentSite(req)

			fieldsToDelete := []string{}

			if site.Logo == image.Id {
				fieldsToDelete = append(fieldsToDelete, "logo")
			}

			if site.Cover == image.Id {
				fieldsToDelete = append(fieldsToDelete, "cover")
			}

			if len(fieldsToDelete) > 0 {
				site.DeleteFields(fieldsToDelete)
			}

			// returns deleted image
			app.render.JSON(rw, http.StatusOK, renderMap{"image": image})
		}
	} else {
		http.NotFound(rw, req)
	}
}

func (app *Application) handleUploadImage(rw http.ResponseWriter, req *http.Request) {
	log.Printf("[handler]: handleUploadImage\n")

	reader, err := req.MultipartReader()
	if err != nil {
		log.Printf("Multipart error: %v", err.Error())
		http.Error(rw, "Failed to parse multipart data", http.StatusBadRequest)
		return
	}

	var fileName string
	var fileContentType string
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

		// @todo Check that content-type is really an image
		fileContentType = part.Header.Get("Content-Type")

		log.Printf("Handling uploaded file: %s", fileName)

		dstPath := utils.AvailableFilePath(path.Join(appUploadDir, fileName))

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
		now := time.Now()

		site := app.getCurrentSite(req)
		if site == nil {
			panic("Site should be set")
		}

		// create image model
		img := &models.Image{
			Id:        bson.NewObjectId(),
			CreatedAt: now,
			UpdatedAt: now,
			SiteId:    site.Id,
			Path:      path.Join(PUBLIC_UPLOAD_PATH, fileInfo.Name()),
			Name:      fileName,
			Size:      fileInfo.Size(),
			Type:      fileContentType,
		}

		if err := app.dbSession.CreateImage(img); err != nil {
			log.Printf("Can't create record: %v - %v", img, err.Error())
			http.Error(rw, "Failed to create image record", http.StatusInternalServerError)
			return
		}

		// @todo Async that the day it becomes problematic
		if err := img.GenerateDerivatives(); err != nil {
			log.Printf("Failed to generate image derivatives: %s - %v", img.Path, err.Error())
		}

		// returns uploaded file path
		app.render.JSON(rw, http.StatusCreated, renderMap{"image": img})
	}
}
