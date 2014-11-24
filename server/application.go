package server

import (
	"github.com/aymerick/kowa/models"
	"github.com/unrolled/render"
)

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
