package server

import (
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
)

// POST /oauth/token
func handleOauthToken(w http.ResponseWriter, req *http.Request) {
	resp := oauthServer.NewResponse()
	defer resp.Close()

	if ar := oauthServer.HandleAccessRequest(resp, req); ar != nil {
		fmt.Println("ar: %v", ar)
		switch ar.Type {
		case osin.PASSWORD:
			// @todo Finish that !
			if ar.Username == "test@test.com" && ar.Password == "test" {
				ar.Authorized = true
			}
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		}
		oauthServer.FinishAccessRequest(resp, req, ar)
	}
	osin.OutputJSON(resp, w, req)
}
