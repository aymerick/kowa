package mailers

import (
	"github.com/spf13/viper"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/models"
)

// Mailer is the interface to all mailers
type Mailer interface {
	Kind() string
	To() string
	Subject() string
}

// BaseMailer is a base for all mailers
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

// NewBaseMailer instanciates a new BaseMailer
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

// Kind is part of Mailer interface
func (mailer *BaseMailer) Kind() string {
	return mailer.kind
}
