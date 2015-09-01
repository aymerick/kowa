package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/helpers"
	"github.com/aymerick/kowa/mailers"
	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/token"
)

type userJson struct {
	User models.User `json:"user"`
}

// POST /api/signup
func (app *Application) handleSignupUser(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	if err := req.ParseForm(); err != nil {
		http.Error(rw, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// get form data
	email := req.Form.Get("email")
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	lang := req.Form.Get("lang")
	tz := req.Form.Get("tz")

	// check lang
	userLang := core.DEFAULT_LANG
	if lang != "" {
		for _, availableLang := range core.Langs {
			if lang == availableLang {
				userLang = lang
				break
			}
		}
	}

	T := i18n.MustTfunc(userLang)

	errors := make(map[string]string)

	// check timezone
	if (tz == "") || !core.ValidTZ(tz) {
		tz = core.DEFAULT_TZ
	}

	// check email format
	emailAddr, err := mail.ParseAddress(email)
	if err != nil || emailAddr.Address == "" {
		errors["email"] = T("signup_email_invalid")
	}

	// check username format
	if username != helpers.NormalizeToUsername(username) {
		errors["username"] = T("signup_username_invalid")
	}

	// check username length
	if len(username) < 4 {
		errors["username"] = T("signup_username_too_short")
	}

	// check password length
	if len(password) < 8 {
		errors["password"] = T("signup_password_too_weak")
	}

	if errors["email"] == "" {
		// check if email is already taken
		if user := currentDBSession.FindUserByEmail(emailAddr.Address); user != nil {
			errors["email"] = T("signup_email_not_available")
		}
	}

	if errors["username"] == "" {
		// check if username is already taken
		if user := currentDBSession.FindUser(username); user != nil {
			errors["username"] = T("signup_username_not_available")
		}
	}

	if len(errors) > 0 {
		app.render.JSON(rw, http.StatusBadRequest, renderMap{"errors": errors})
		return
	}

	// encrypt password
	encryptedPassword, errPass := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errPass != nil {
		http.Error(rw, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	// insert user
	user := &models.User{
		Id:       username,
		Email:    emailAddr.Address,
		Status:   models.USER_STATUS_PENDING,
		Lang:     userLang,
		Password: string(encryptedPassword),
	}

	if err := currentDBSession.CreateUser(user); err != nil {
		http.Error(rw, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// send signup confirmation email
	go mailers.NewSignupMailer(user).Send()

	app.render.JSON(rw, http.StatusCreated, renderMap{"user": user})
}

func (app *Application) userFromSignupValidationToken(rw http.ResponseWriter, req *http.Request, checkExpiration bool) *models.User {
	currentDBSession := app.getCurrentDBSession(req)

	if err := req.ParseForm(); err != nil {
		http.Error(rw, "Failed to parse form data", http.StatusBadRequest)
		return nil
	}

	tokenStr := req.Form.Get("token")
	if tokenStr == "" {
		// Missing token
		unauthorized(rw)
		return nil
	}

	tok := token.Decode(tokenStr)
	if tok == nil {
		// Invalid token
		unauthorized(rw)
		return nil
	}

	userId := tok.AccountValidationUser()
	if userId == "" {
		// Erroneous token
		unauthorized(rw)
		return nil
	}

	if checkExpiration && tok.Expired() {
		http.Error(rw, "Signup validation token expired", http.StatusUnauthorized)
		return nil
	}

	// fetch user
	user := currentDBSession.FindUser(userId)
	if user == nil {
		http.NotFound(rw, req)
		return nil
	}

	return user
}

// POST /api/signup/validate
func (app *Application) handleSignupValidate(rw http.ResponseWriter, req *http.Request) {
	// get user
	user := app.userFromSignupValidationToken(rw, req, true)
	if user == nil {
		// there was an error
		return
	}

	// check if user was already validated
	if !user.AccountValidated() {
		// mark user as validated
		user.SetAccountValidated()
	}

	// @todo Auto login ?

	// @todo Send Welcome email

	app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
}

// POST /api/signup/sendmail
func (app *Application) handleSignupSendMail(rw http.ResponseWriter, req *http.Request) {
	// get user without checking token expiration
	user := app.userFromSignupValidationToken(rw, req, false)
	if user == nil {
		// there was an error
		return
	}

	// send signup confirmation email
	go mailers.NewSignupMailer(user).Send()

	app.render.JSON(rw, http.StatusOK, renderMap{"response": "ok"})
}

// GET /api/me
func (app *Application) handleGetMe(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	currentUser := app.getCurrentUser(req)
	userId := currentUser.Id

	if user := currentDBSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}
func (app *Application) handleGetUser(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	vars := mux.Vars(req)
	userId := vars["user_id"]

	if user := currentDBSession.FindUser(userId); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /api/users/{user_id}
func (app *Application) handleUpdateUser(rw http.ResponseWriter, req *http.Request) {
	user := app.getCurrentUser(req)
	if user != nil {
		var reqJson userJson

		if err := json.NewDecoder(req.Body).Decode(&reqJson); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		// @todo [security] Check all fields !
		_, err := user.Update(&reqJson.User)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to update user", http.StatusInternalServerError)
			return
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}/sites
func (app *Application) handleGetUserSites(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	vars := mux.Vars(req)
	userId := vars["user_id"]

	// check current user
	currentUser := app.getCurrentUser(req)
	if currentUser == nil {
		unauthorized(rw)
		return
	}

	if currentUser.Id != userId {
		unauthorized(rw)
		return
	}

	if user := currentDBSession.FindUser(userId); user != nil {
		var image *models.Image
		images := []*models.Image{}

		pageSettingsArray := []*models.SitePageSettings{}

		sites := user.FindSites()
		for _, site := range *sites {
			if image = site.FindLogo(); image != nil {
				images = append(images, image)
			}

			if image = site.FindCover(); image != nil {
				images = append(images, image)
			}

			for _, pageSettings := range site.PageSettings {
				pageSettingsArray = append(pageSettingsArray, pageSettings)

				if image = site.FindPageSettingsCover(pageSettings.Kind); image != nil {
					images = append(images, image)
				}
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{"sites": sites, "images": images, "sitePageSettings": pageSettingsArray})
	} else {
		http.NotFound(rw, req)
	}
}
