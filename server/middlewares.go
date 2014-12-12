package server

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/RangelReale/osin"
	"github.com/aymerick/kowa/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2/bson"
)

// middleware: logs requests
func (app *Application) loggingMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		startAt := time.Now()
		next.ServeHTTP(rw, req)
		endAt := time.Now()

		log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), endAt.Sub(startAt))
	}

	return http.HandlerFunc(fn)
}

// middleware: recovers panic
func (app *Application) recoveryMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)

				stack := debug.Stack()
				log.Printf("PANIC: %s\n%s", err, stack)
			}
		}()

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: injects CORS headers
func (app *Application) corsMiddleware() func(next http.Handler) http.Handler {
	result := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization"},
		// AllowCredentials: true,
		// MaxAge:           5 * time.Minute,
	})

	return result.Handler
}

// middleware: ensures user is authenticated and injects 'currentUser' in context
func (app *Application) ensureAuthMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		var err error

		// ex: Authorization: Bearer Zjg5ZmEwNDYtNGI3NS00MTk4LWFhYzgtZmVlNGRkZDQ3YzAx
		authValue := req.Header.Get("Authorization")
		if len(authValue) < 7 || authValue[:7] != "Bearer " {
			unauthorized(rw)
			return
		}

		var accessData *osin.AccessData
		accessData, err = app.oauthServer.Storage.LoadAccess(authValue[7:])
		if err != nil {
			unauthorized(rw)
			return
		}

		// @todo Check accessData.CreatedAt

		userId, ok := accessData.UserData.(string)
		if !ok || userId == "" {
			unauthorized(rw)
			return
		}

		if currentUser := app.dbSession.FindUser(userId); currentUser != nil {
			log.Printf("Current user is: %s [%s]\n", currentUser.Fullname(), userId)
			context.Set(req, "currentUser", currentUser)
		} else {
			unauthorized(rw)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures that currently authenticated user is allowed to access a /users/{user_id}/* requests
func (app *Application) ensureUserAccessMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		// check current user
		currentUser := app.getCurrentUser(req)
		if currentUser == nil {
			unauthorized(rw)
			return
		}

		vars := mux.Vars(req)
		userId := vars["user_id"]

		// check that current user only access his stuff
		if userId != currentUser.Id {
			unauthorized(rw)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures site exists and injects 'currentSite' in context
func (app *Application) ensureSiteMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		var currentSite *models.Site

		// site id
		siteId := vars["site_id"]
		if siteId != "" {
			currentSite = app.dbSession.FindSite(siteId)
		} else {
			// post
			currentPost := app.getCurrentPost(req)
			if currentPost != nil {
				currentSite = currentPost.FindSite()
			} else {
				// image
				currentImage := app.getCurrentImage(req)
				if currentImage != nil {
					currentSite = currentImage.FindSite()
				}
			}
		}

		if currentSite != nil {
			log.Printf("Current site is: %s [%s]\n", currentSite.Name, siteId)
			context.Set(req, "currentSite", currentSite)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures that currently authenticated user is allowed to access a /sites/{site_id}/* requests
func (app *Application) ensureSiteOwnerAccessMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		// check current user
		currentUser := app.getCurrentUser(req)
		if currentUser == nil {
			panic("Should be auth")
		}

		// check current site
		currentSite := app.getCurrentSite(req)
		if currentSite == nil {
			panic("Should have site")
		}

		if currentSite.UserId != currentUser.Id {
			unauthorized(rw)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures post exists and injects 'currentPost' in context
func (app *Application) ensurePostMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		postId := vars["post_id"]
		if postId == "" {
			panic("Should have post_id")
		}

		if currentPost := app.dbSession.FindPost(bson.ObjectIdHex(postId)); currentPost != nil {
			context.Set(req, "currentPost", currentPost)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures image exists and injects 'currentImage' in context
func (app *Application) ensureImageMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		imageId := vars["image_id"]
		if imageId == "" {
			panic("Should have image_id")
		}

		if currentImage := app.dbSession.FindImage(bson.ObjectIdHex(imageId)); currentImage != nil {
			context.Set(req, "currentImage", currentImage)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}
