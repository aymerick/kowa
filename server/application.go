package server

import (
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/gorilla/context"
	"github.com/spf13/viper"
	"github.com/unrolled/render"

	"github.com/aymerick/kowa/models"
)

type Application struct {
	port        string
	render      *render.Render
	dbSession   *models.DBSession
	oauthServer *osin.Server
}

func NewApplication() *Application {
	dbSession := models.NewDBSession()
	dbSession.EnsureIndexes()

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
		dbSession:   dbSession,
		oauthServer: oauthServer,
	}
}

//
// Request context
//

func (app *Application) getCurrentUser(req *http.Request) *models.User {
	if currentUser := context.Get(req, "currentUser"); currentUser != nil {
		return currentUser.(*models.User)
	}
	return nil
}

func (app *Application) getCurrentSite(req *http.Request) *models.Site {
	if currentSite := context.Get(req, "currentSite"); currentSite != nil {
		return currentSite.(*models.Site)
	}
	return nil
}

func (app *Application) getCurrentPost(req *http.Request) *models.Post {
	if currentPost := context.Get(req, "currentPost"); currentPost != nil {
		return currentPost.(*models.Post)
	}
	return nil
}

func (app *Application) getCurrentImage(req *http.Request) *models.Image {
	if currentImage := context.Get(req, "currentImage"); currentImage != nil {
		return currentImage.(*models.Image)
	}
	return nil
}
