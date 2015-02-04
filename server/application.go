package server

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

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
	buildMaster *BuildMaster
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
		buildMaster: NewBuildMaster(),
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
		go app.serveBuiltSites()
	}

	// start web server
	log.Println("Running on port:", app.port)
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

// Called when some content changed on given site
func (app *Application) onSiteChange(site *models.Site) {
	// update ChangedAt anchor
	site.SetChangedAt(time.Now())

	// rebuild changed site
	app.buildSite(site)
}

// DEBUG: Serve built sites
func (app *Application) serveBuiltSites() {
	dir := path.Join(viper.GetString("working_dir"), viper.GetString("output_dir"))
	port := viper.GetInt("serve_output_port")

	log.Println("Serving built sites on port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), http.FileServer(http.Dir(dir))))
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

func (app *Application) getCurrentImage(req *http.Request) *models.Image {
	if currentImage := context.Get(req, "currentImage"); currentImage != nil {
		return currentImage.(*models.Image)
	}
	return nil
}
