package server

import (
	"net/http"

	"github.com/RangelReale/osin"
)

type OauthController struct {
	*ApplicationController
}

// POST /oauth/token
func (this *OauthController) handleOauthToken(rw http.ResponseWriter, req *http.Request) error {
	resp := oauthServer.NewResponse()
	defer resp.Close()

	if ar := oauthServer.HandleAccessRequest(resp, req); ar != nil {
		switch ar.Type {
		case osin.PASSWORD:
			// @todo Finish that !
			if ar.Username == "test@test.com" && ar.Password == "test" {
				ar.UserData = "test@test.com"
				ar.Authorized = true
			}
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		}
		oauthServer.FinishAccessRequest(resp, req, ar)
	}

	osin.OutputJSON(resp, rw, req)

	return nil
}

// POST /oauth/revoke
func (this *OauthController) handleOauthRevoke(rw http.ResponseWriter, req *http.Request) error {
	resp := oauthServer.NewResponse()
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

	return nil
}
