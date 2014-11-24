package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
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
