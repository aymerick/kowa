package server

import (
	"net/http"

	"github.com/aymerick/kowa/models"
	"github.com/unrolled/render"
)

type Action func(rw http.ResponseWriter, r *http.Request) error

type Application struct {
	render    *render.Render
	dbSession *models.DBSession
}

func NewApplication() *Application {
	return &Application{
		render:    render.New(render.Options{}),
		dbSession: models.NewDBSession(),
	}
}
