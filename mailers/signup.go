package mailers

import (
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
)

// Implements Mailer
type SignupMailer struct {
	user *models.User
	T    i18n.TranslateFunc
}

func NewSignupMailer(user *models.User) *SignupMailer {
	return &SignupMailer{
		user: user,
		T:    i18n.MustTfunc(user.Lang),
	}
}

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
	return mailer.T("Please confirm your signup")
}

func (mailer *SignupMailer) Html() string {
	return "@todo HTML content"
}

func (mailer *SignupMailer) Text() string {
	return "@todo Text content"
}
