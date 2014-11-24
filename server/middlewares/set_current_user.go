package middlewares

import "net/http"

func SetCurrentUser(next http.Handler) http.Handler {
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
