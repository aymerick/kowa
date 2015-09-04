package server

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/models"
)

type upload struct {
	name  string
	ctype string
	info  os.FileInfo
}

func newUpload() *upload {
	return &upload{}
}

func handleUpload(rw http.ResponseWriter, req *http.Request, site *models.Site, allowedTypes []string) *upload {
	reader, err := req.MultipartReader()
	if err != nil {
		log.Printf("Multipart error: %v", err.Error())
		http.Error(rw, "Failed to parse multipart data", http.StatusBadRequest)
		return nil
	}

	result := newUpload()

	for result.name == "" {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		result.name = part.FileName()
		if result.name == "" {
			continue
		}

		// Check content type
		result.ctype = part.Header.Get("Content-Type")

		if !allowedContentType(result.ctype, allowedTypes) {
			log.Printf("Unsupported content type for file upload: %v", result.ctype)
			http.Error(rw, "Unsupported content type", http.StatusBadRequest)
			return nil
		}

		// copy uploaded file
		log.Printf("Handling uploaded file: %s", result.name)

		dstPath := helpers.AvailableFilePath(core.UploadSiteFilePath(site.Id, result.name))

		dst, err := os.Create(dstPath)
		if err != nil {
			log.Printf("Can't create file: %s - %v", dstPath, err.Error())
			http.Error(rw, "Failed to create uploaded file", http.StatusInternalServerError)
			return nil
		}

		defer dst.Close()

		if _, err := io.Copy(dst, part); err != nil {
			log.Printf("Can't save file: %s - %v", dstPath, err.Error())
			http.Error(rw, "Failed to save uploaded file", http.StatusInternalServerError)
			return nil
		}

		var errStat error
		result.info, errStat = os.Stat(dstPath)
		if os.IsNotExist(errStat) {
			http.Error(rw, "Failed to create uploaded file", http.StatusInternalServerError)
			return nil
		}
	}

	if result.name == "" {
		http.Error(rw, "File not found in multipart", http.StatusBadRequest)
		return nil
	}

	return result
}

func allowedContentType(ct string, allowed []string) bool {
	for _, allowedCT := range allowed {
		if ct == allowedCT {
			return true
		}
	}

	return false
}
