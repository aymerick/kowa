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
		log.Printf("[middleware]: loggingMiddleware\n")

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
		log.Printf("[middleware]: recoveryMiddleware\n")

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
		log.Printf("[middleware]: ensureAuthMiddleware\n")

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

// middleware: ensures site exists and injects 'currentSite' in context
func (app *Application) ensureSiteMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("[middleware]: ensureSiteMiddleware\n")

		vars := mux.Vars(req)
		siteId := vars["site_id"]

		if currentSite := app.dbSession.FindSite(bson.ObjectIdHex(siteId)); currentSite != nil {
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

// middleware: ensures that currently authenticated user is allowed to access a /users/{user_id}/* requests
func (app *Application) ensureUserAccessMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("[middleware]: ensureUserAccessMiddleware\n")

		// check current user
		currentUser := context.Get(req, "currentUser").(*models.User)
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

// middleware: ensures that currently authenticated user is allowed to access a /sites/{site_id}/* requests
func (app *Application) ensureSiteOwnerAccessMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("[middleware]: ensureSiteOwnerAccessMiddleware\n")

		// check current user
		currentUser := context.Get(req, "currentUser").(*models.User)
		if currentUser == nil {
			panic("Should be auth")
		}

		// check current site
		currentSite := context.Get(req, "currentSite").(*models.Site)
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
