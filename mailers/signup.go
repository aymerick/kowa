package mailers

import (
	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/token"
)

// Implements Mailer
type SignupMailer struct {
	*BaseMailer

	// Template variables
	Email         string
	ActivationUrl string
}

func NewSignupMailer(user *models.User) *SignupMailer {
	result := &SignupMailer{
		BaseMailer: NewBaseMailer("signup", user),

		// Template variables
		Email:         user.Email,
		ActivationUrl: token.AccountActivationUrl(user),
	}

	result.I18n = result.computeI18n()

	return result
}

// Send mail
func (mailer *SignupMailer) Send() error {
	return NewSender(mailer).Send()
}

// Computes translations
func (mailer *SignupMailer) computeI18n() map[string]string {
	return map[string]string{
		"thanks":                mailer.T("signup_email_thanks", core.P{"ServiceName": mailer.ServiceName}),
		"one_more_step":         mailer.T("signup_email_one_more_step"),
		"activate_your_account": mailer.T("signup_email_activate_your_account", core.P{"ActivationUrl": mailer.ActivationUrl}),
		"click_button":          mailer.T("signup_email_click_button"),
		"activate_account":      mailer.T("signup_email_activate_account"),
	}
}

//
// Mailer interface
//

func (mailer *SignupMailer) To() string {
	return mailer.user.MailAddress()
}

func (mailer *SignupMailer) Subject() string {
	return mailer.T("signup_email_subject", core.P{"ServiceName": mailer.ServiceName})
}
