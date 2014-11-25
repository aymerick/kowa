package server

import (
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/RangelReale/osin"

	"github.com/aymerick/kowa/models"
)

// POST /oauth/token
func (app *Application) handleOauthToken(rw http.ResponseWriter, req *http.Request) {
	resp := app.oauthServer.NewResponse()
	defer resp.Close()

	if ar := app.oauthServer.HandleAccessRequest(resp, req); ar != nil {
		switch ar.Type {
		case osin.PASSWORD:
			var user *models.User

			user = app.dbSession.FindUser(ar.Username)
			if user == nil {
				user = app.dbSession.FindUserByEmail(ar.Username)
			}

			if user != nil {
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ar.Password))
				if err == nil {
					ar.UserData = user.Id
					ar.Authorized = true
				}
			}
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		}
		app.oauthServer.FinishAccessRequest(resp, req, ar)
	}

	osin.OutputJSON(resp, rw, req)
}

// POST /oauth/revoke
func (app *Application) handleOauthRevoke(rw http.ResponseWriter, req *http.Request) {
	resp := app.oauthServer.NewResponse()
	defer resp.Close()

	err := req.ParseForm()
	if err != nil {
		resp.SetError(osin.E_INVALID_REQUEST, "")
		resp.InternalError = err
	} else {
		tokenType := req.Form.Get("token_type_hint")
		token := req.Form.Get("token")

		switch tokenType {
		case "access_token":
			resp.Storage.RemoveAccess(token)

		case "refresh_token":
			resp.Storage.RemoveRefresh(token)
		}
	}

	osin.OutputJSON(resp, rw, req)
}
