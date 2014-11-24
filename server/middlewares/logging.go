package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		startAt := time.Now()
		next.ServeHTTP(rw, r)
		endAt := time.Now()

		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), endAt.Sub(startAt))
	}

	return http.HandlerFunc(fn)
}
