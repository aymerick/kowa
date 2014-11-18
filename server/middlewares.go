package server

import (
	"github.com/codegangsta/negroni"
	cors "github.com/zimmski/negroni-cors"
)

func setupMiddlewares(n *negroni.Negroni) {
	n.Use(cors.NewAllow(&cors.Options{
		AllowAllOrigins: true,
		// AllowCredentials: true,
		AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization"},
		// MaxAge:           5 * time.Minute,
	}))
}
