package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/aymerick/kowa/utils"
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

	var partFileName string
	var fileInfo os.FileInfo

	for partFileName == "" {
		log.Printf("partFileName: %s -> Next Part()", partFileName)

		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		partFileName = part.FileName()
		if partFileName == "" {
			continue
		}

		log.Printf("Handling uploaded file: %s", partFileName)

		dstPath := utils.AvailableFilePath(path.Join(appUploadDir, partFileName))

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

	if partFileName == "" {
		http.Error(rw, "File data not found in multipart", http.StatusBadRequest)
	} else {
		// returns uploaded file path
		app.render.JSON(rw, http.StatusCreated, renderMap{"name": partFileName, "path": path.Join(PUBLIC_UPLOAD_PATH, fileInfo.Name()), "size": fileInfo.Size()})
	}
}
