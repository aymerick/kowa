package server

import (
	"github.com/RangelReale/osin"
	"github.com/aymerick/kowa/models"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

type Application struct {
	port        string
	render      *render.Render
	dbSession   *models.DBSession
	oauthServer *osin.Server
}

func NewApplication() *Application {
	// setup osin oauth2 server
	osinConfig := osin.NewServerConfig()
	osinConfig.AccessExpiration = 3600 // One hour
	osinConfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.TOKEN}
	osinConfig.AllowedAccessTypes = osin.AllowedAccessType{osin.PASSWORD, osin.REFRESH_TOKEN}
	osinConfig.ErrorStatusCode = 401

	oauthServer := osin.NewServer(osinConfig, NewOAuthStorage())

	return &Application{
		port:        viper.GetString("port"),
		render:      render.New(render.Options{}),
		dbSession:   models.NewDBSession(),
		oauthServer: oauthServer,
	}
}
