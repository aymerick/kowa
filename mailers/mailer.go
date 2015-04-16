package mailers

import (
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/models"
)

// provides mail data to Sender
type Mailer interface {
	Kind() string
	To() string
	Subject() string
}

type BaseMailer struct {
	kind string
	user *models.User

	T core.TranslateFunc

	// Template variables
	I18n                   map[string]string
	ServiceName            string
	ServiceLogo            string
	ServiceUrl             string
	ServicePostalAddress   string
	ServiceCopyrightNotice string
}

func NewBaseMailer(kind string, user *models.User) *BaseMailer {
	return &BaseMailer{
		kind: kind,
		user: user,

		T: core.MustTfunc(user.Lang),

		// Template variables
		ServiceName:            viper.GetString("service_name"),
		ServiceLogo:            viper.GetString("service_logo"),
		ServiceUrl:             viper.GetString("service_url"),
		ServicePostalAddress:   viper.GetString("service_postal_address"),
		ServiceCopyrightNotice: viper.GetString("service_copyright_notice"),
	}
}

//
// Mailer interface
//

func (mailer *BaseMailer) Kind() string {
	return mailer.kind
}
