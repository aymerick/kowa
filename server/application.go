package server

import (
	"net/http"

	"github.com/unrolled/render"
)

type Action func(rw http.ResponseWriter, r *http.Request) error

// Application Controller
type ApplicationController struct {
	render *render.Render
}

func (this *ApplicationController) Action(action Action) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := action(rw, r); err != nil {
			http.Error(rw, err.Error(), 500)
		}
	})
}
