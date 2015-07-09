package server

import (
	"log"
	"net/http"
	"time"

	"github.com/RangelReale/osin"
	"github.com/gorilla/context"
	"github.com/spf13/viper"
	"github.com/unrolled/render"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/mailers"
	"github.com/aymerick/kowa/models"
)

type Application struct {
	port         string
	render       *render.Render
	dbSession    *models.DBSession
	oauthStorage *OAuthStorage
	oauthServer  *osin.Server
	buildMaster  *BuildMaster
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

	oauthStorage := NewOAuthStorage()
	oauthServer := osin.NewServer(osinConfig, oauthStorage)

	return &Application{
		port:         viper.GetString("port"),
		render:       render.New(render.Options{}),
		dbSession:    dbSession,
		oauthStorage: oauthStorage,
		oauthServer:  oauthServer,
		buildMaster:  NewBuildMaster(),
	}
}

// Setup application server
func (app *Application) Setup() {
	core.EnsureUploadDir()

	// Ensure oauth client
	if err := app.oauthStorage.EnsureOAuthClient(); err != nil {
		panic(err)
	}

	if viper.GetString("mail_tpl_dir") != "" {
		// Set templates dir for mails
		mailers.SetTemplatesDir(viper.GetString("mail_tpl_dir"))
	}
}

// Run application server
func (app *Application) Run() {
	// @todo Build pending sites on startup ?! (The ones with BuiltAt < UpdatedAt)
	//       And add a command "kowa build_pending" too

	// start build master
	app.buildMaster.run()

	// TODO: only for dev
	if viper.GetBool("serve_output") {
		go app.buildMaster.serveSites()
	}

	// start web server
	log.Println("Running Web Server on port:", app.port)
	http.ListenAndServe(":"+app.port, app.newWebRouter())
}

// Stop application server
func (app *Application) Stop() {
	// stop build master
	app.buildMaster.stop()
}

// Build site
func (app *Application) buildSite(site *models.Site) {
	app.buildMaster.launchSiteBuild(site)
}

// Delete built site
func (app *Application) deleteBuild(site *models.Site) {
	app.buildMaster.launchSiteDeletion(site)
}

// Called when some content changed on given site
func (app *Application) onSiteChange(site *models.Site) {
	// update ChangedAt anchor
	site.SetChangedAt(time.Now())

	// rebuild changed site
	app.buildSite(site)
}

// Called when site is deleted
func (app *Application) onSiteDeletion(site *models.Site) {
	// delete build
	app.deleteBuild(site)
}

//
// Request context
//

func (app *Application) getCurrentDBSession(req *http.Request) *models.DBSession {
	if currentDBSession := context.Get(req, "currentDBSession"); currentDBSession != nil {
		return currentDBSession.(*models.DBSession)
	}
	return nil
}

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

func (app *Application) getCurrentEvent(req *http.Request) *models.Event {
	if currentEvent := context.Get(req, "currentEvent"); currentEvent != nil {
		return currentEvent.(*models.Event)
	}
	return nil
}

func (app *Application) getCurrentPage(req *http.Request) *models.Page {
	if currentPage := context.Get(req, "currentPage"); currentPage != nil {
		return currentPage.(*models.Page)
	}
	return nil
}

func (app *Application) getCurrentActivity(req *http.Request) *models.Activity {
	if currentActivity := context.Get(req, "currentActivity"); currentActivity != nil {
		return currentActivity.(*models.Activity)
	}
	return nil
}

func (app *Application) getCurrentMember(req *http.Request) *models.Member {
	if currentMember := context.Get(req, "currentMember"); currentMember != nil {
		return currentMember.(*models.Member)
	}
	return nil
}

func (app *Application) getCurrentImage(req *http.Request) *models.Image {
	if currentImage := context.Get(req, "currentImage"); currentImage != nil {
		return currentImage.(*models.Image)
	}
	return nil
}

func (app *Application) getCurrentFile(req *http.Request) *models.File {
	if currentFile := context.Get(req, "currentFile"); currentFile != nil {
		return currentFile.(*models.File)
	}
	return nil
}

//
// Endpoints
//

// GET /api/configuration
func (app *Application) handleGetConfig(rw http.ResponseWriter, req *http.Request) {
	result := renderMap{
		"langs": []map[string]string{
			{"id": "en", "name": "English"},
			{"id": "fr", "name": "Français"},
		},
		"formats": []map[string]string{
			{"id": "html", "name": "Rich Text"},
			{"id": "md", "name": "Markdown"},
		},
		"themes":  []string{"ailes", "willy"},
		"domains": viper.GetStringSlice("service_domains"),
	}

	app.render.JSON(rw, http.StatusOK, result)
}
