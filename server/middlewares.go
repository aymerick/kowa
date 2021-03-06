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

// middleware: setup database session
func (app *Application) dbSessionMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		dbSessionCopy := app.dbSession.Copy()

		context.Set(req, "currentDBSession", dbSessionCopy)
		defer dbSessionCopy.Close()

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

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

// middleware: ensures user is NOT authenticated
func (app *Application) ensureNotAuthMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		// ex: Authorization: Bearer Zjg5ZmEwNDYtNGI3NS00MTk4LWFhYzgtZmVlNGRkZDQ3YzAx
		authValue := req.Header.Get("Authorization")
		if len(authValue) > 0 {
			unauthorized(rw)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
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

		userID, ok := accessData.UserData.(string)
		if !ok || userID == "" {
			unauthorized(rw)
			return
		}

		currentDBSession := app.getCurrentDBSession(req)

		if currentUser := currentDBSession.FindUser(userID); currentUser != nil {
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
		userID := vars["user_id"]

		// check that current user only access his stuff
		if userID != currentUser.ID {
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
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)

		var currentSite *models.Site

		// site id
		siteID := vars["site_id"]
		if siteID != "" {
			currentSite = currentDBSession.FindSite(siteID)
		}

		// post
		if currentSite == nil {
			currentPost := app.getCurrentPost(req)
			if currentPost != nil {
				currentSite = currentPost.FindSite()
			}
		}

		// event
		if currentSite == nil {
			currentEvent := app.getCurrentEvent(req)
			if currentEvent != nil {
				currentSite = currentEvent.FindSite()
			}
		}

		// page
		if currentSite == nil {
			currentPage := app.getCurrentPage(req)
			if currentPage != nil {
				currentSite = currentPage.FindSite()
			}
		}

		// activity
		if currentSite == nil {
			currentActivity := app.getCurrentActivity(req)
			if currentActivity != nil {
				currentSite = currentActivity.FindSite()
			}
		}

		// member
		if currentSite == nil {
			currentMember := app.getCurrentMember(req)
			if currentMember != nil {
				currentSite = currentMember.FindSite()
			}
		}

		// image
		if currentSite == nil {
			currentImage := app.getCurrentImage(req)
			if currentImage != nil {
				currentSite = currentImage.FindSite()
			}
		}

		// file
		if currentSite == nil {
			currentFile := app.getCurrentFile(req)
			if currentFile != nil {
				currentSite = currentFile.FindSite()
			}
		}

		if currentSite != nil {
			// log.Printf("Current site is: %s [%s]\n", currentSite.Name, siteID)
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

		if currentSite.UserID != currentUser.ID {
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
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)
		postID := vars["post_id"]
		if postID == "" {
			panic("Should have post_id")
		}

		if currentPost := currentDBSession.FindPost(bson.ObjectIdHex(postID)); currentPost != nil {
			context.Set(req, "currentPost", currentPost)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures event exists and injects 'currentEvent' in context
func (app *Application) ensureEventMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)
		eventID := vars["event_id"]
		if eventID == "" {
			panic("Should have event_id")
		}

		if currentEvent := currentDBSession.FindEvent(bson.ObjectIdHex(eventID)); currentEvent != nil {
			context.Set(req, "currentEvent", currentEvent)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures page exists and injects 'currentPage' in context
func (app *Application) ensurePageMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)
		pageID := vars["page_id"]
		if pageID == "" {
			panic("Should have page_id")
		}

		if currentPage := currentDBSession.FindPage(bson.ObjectIdHex(pageID)); currentPage != nil {
			context.Set(req, "currentPage", currentPage)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures activity exists and injects 'currentActivity' in context
func (app *Application) ensureActivityMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)
		activityID := vars["activity_id"]
		if activityID == "" {
			panic("Should have activity_id")
		}

		if currentActivity := currentDBSession.FindActivity(bson.ObjectIdHex(activityID)); currentActivity != nil {
			context.Set(req, "currentActivity", currentActivity)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures member exists and injects 'currentMember' in context
func (app *Application) ensureMemberMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)
		memberID := vars["member_id"]
		if memberID == "" {
			panic("Should have member_id")
		}

		if currentMember := currentDBSession.FindMember(bson.ObjectIdHex(memberID)); currentMember != nil {
			context.Set(req, "currentMember", currentMember)
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
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)
		imageID := vars["image_id"]
		if imageID == "" {
			panic("Should have image_id")
		}

		if currentImage := currentDBSession.FindImage(bson.ObjectIdHex(imageID)); currentImage != nil {
			context.Set(req, "currentImage", currentImage)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// middleware: ensures file exists and injects 'currentFile' in context
func (app *Application) ensureFileMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		currentDBSession := app.getCurrentDBSession(req)

		vars := mux.Vars(req)
		fileID := vars["file_id"]
		if fileID == "" {
			panic("Should have file_id")
		}

		if currentFile := currentDBSession.FindFile(bson.ObjectIdHex(fileID)); currentFile != nil {
			context.Set(req, "currentFile", currentFile)
		} else {
			http.NotFound(rw, req)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}
