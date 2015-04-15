package mailers

import (
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
)

// Implements Mailer
type SignupMailer struct {
	*BaseMailer

	user *models.User
	T    i18n.TranslateFunc

	// Template variables
	Username       string
	Email          string
	ActivationLink string
}

func NewSignupMailer(user *models.User) *SignupMailer {
	return &SignupMailer{
		BaseMailer: NewBaseMailer("signup"),

		user: user,
		T:    i18n.MustTfunc(user.Lang),

		// Template variables
		Username:       user.Id,
		Email:          user.Email,
		ActivationLink: user.ActivationLink(),
	}
}

// Send mail
func (mailer *SignupMailer) Send() error {
	return NewSender(mailer).Send()
}

//
// Mailer interface
//

func (mailer *SignupMailer) To() string {
	return mailer.user.MailAddress()
}

func (mailer *SignupMailer) Subject() string {
	return mailer.T("signup_email_subject", map[string]interface{}{"ServiceName": mailer.ServiceName})
}
