package server

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

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

func (app *Application) setCurrentUserMiddleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		// ex: Authorization: Bearer Zjg5ZmEwNDYtNGI3NS00MTk4LWFhYzgtZmVlNGRkZDQ3YzAx
		authValue := r.Header.Get("Authorization")
		if len(authValue) < 7 || authValue[:7] != "Bearer " {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}

		// authToken := authValue[7:]

		// LoadAccess()

		// log.Printf("Current user is: %s\n", currentUser)

		next.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(fn)
}
