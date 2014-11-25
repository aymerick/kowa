package server

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/RangelReale/osin"
	"github.com/gorilla/context"
	"github.com/rs/cors"
)

func (app *Application) loggingMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		startAt := time.Now()
		next.ServeHTTP(rw, r)
		endAt := time.Now()

		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), endAt.Sub(startAt))
	}

	return http.HandlerFunc(fn)
}

func (app *Application) recoveryMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)

				stack := debug.Stack()
				log.Printf("PANIC: %s\n%s", err, stack)
			}
		}()

		next.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(fn)
}

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

func (app *Application) ensureAuthMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		var err error

		// ex: Authorization: Bearer Zjg5ZmEwNDYtNGI3NS00MTk4LWFhYzgtZmVlNGRkZDQ3YzAx
		authValue := r.Header.Get("Authorization")
		if len(authValue) < 7 || authValue[:7] != "Bearer " {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}

		var accessData *osin.AccessData
		accessData, err = app.oauthServer.Storage.LoadAccess(authValue[7:])
		if err != nil {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}

		// @todo Check accessData.CreatedAt

		userId, ok := accessData.UserData.(string)
		if !ok || userId == "" {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}

		if currentUser := app.dbSession.FindUser(userId); currentUser != nil {
			log.Printf("Current user is: %s [%s]\n", currentUser.Fullname(), userId)
			context.Set(r, "currentUser", currentUser)
		} else {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(fn)
}
