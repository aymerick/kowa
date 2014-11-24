package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

func Cors() func(next http.Handler) http.Handler {
	result := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization"},
		// AllowCredentials: true,
		// MaxAge:           5 * time.Minute,
	})

	return result.Handler
}
