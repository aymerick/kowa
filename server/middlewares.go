package server

import (
	"net/http"

	"github.com/codegangsta/negroni"
	cors "github.com/zimmski/negroni-cors"
)

func setupCORSMiddleware(n *negroni.Negroni) {
	// CORS
	n.Use(cors.NewAllow(&cors.Options{
		AllowAllOrigins: true,
		// AllowCredentials: true,
		AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization"},
		// MaxAge:           5 * time.Minute,
	}))
}

// Injects an Authorization heaer wih a fake client_id / client_secret pair, to make osin lib happy
//
// alternative: https://github.com/simplabs/ember-simple-auth/issues/226
func injectOauthSecretMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// set fake secret
	r.Header.Set("Authorization", "Basic "+CLIENT_AUTH_VALUE)

	// do some stuff before
	next(rw, r)
}
