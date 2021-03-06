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

type userJSON struct {
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
	userLang := core.DefaultLang
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
		tz = core.DefaultTZ
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
		ID:       username,
		Email:    emailAddr.Address,
		Status:   models.UserStatusPending,
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

	userID := tok.AccountValidationUser()
	if userID == "" {
		// Erroneous token
		unauthorized(rw)
		return nil
	}

	if checkExpiration && tok.Expired() {
		http.Error(rw, "Signup validation token expired", http.StatusUnauthorized)
		return nil
	}

	// fetch user
	user := currentDBSession.FindUser(userID)
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
	userID := currentUser.ID

	if user := currentDBSession.FindUser(userID); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// GET /api/users/{user_id}
func (app *Application) handleGetUser(rw http.ResponseWriter, req *http.Request) {
	currentDBSession := app.getCurrentDBSession(req)

	vars := mux.Vars(req)
	userID := vars["user_id"]

	if user := currentDBSession.FindUser(userID); user != nil {
		app.render.JSON(rw, http.StatusOK, renderMap{"user": user})
	} else {
		http.NotFound(rw, req)
	}
}

// PUT /api/users/{user_id}
func (app *Application) handleUpdateUser(rw http.ResponseWriter, req *http.Request) {
	user := app.getCurrentUser(req)
	if user != nil {
		var reqJSON userJSON

		if err := json.NewDecoder(req.Body).Decode(&reqJSON); err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(rw, "Failed to decode JSON data", http.StatusBadRequest)
			return
		}

		// @todo [security] Check all fields !
		_, err := user.Update(&reqJSON.User)
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
	userID := vars["user_id"]

	// check current user
	currentUser := app.getCurrentUser(req)
	if currentUser == nil {
		unauthorized(rw)
		return
	}

	if currentUser.ID != userID {
		unauthorized(rw)
		return
	}

	if user := currentDBSession.FindUser(userID); user != nil {
		var image *models.Image
		images := []*models.Image{}

		pageSettingsArray := []*models.SitePageSettings{}
		themeSettingsArray := []*models.SiteThemeSettings{}

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

			for themeID, themeSettings := range site.ThemeSettings {
				themeSettings.Theme = themeID

				themeSettingsArray = append(themeSettingsArray, themeSettings)
			}
		}

		app.render.JSON(rw, http.StatusOK, renderMap{
			"sites":             sites,
			"images":            images,
			"sitePageSettings":  pageSettingsArray,
			"siteThemeSettings": themeSettingsArray,
		})
	} else {
		http.NotFound(rw, req)
	}
}
