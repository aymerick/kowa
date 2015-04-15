package mailers

import "github.com/spf13/viper"

// provides mail data to Sender
type Mailer interface {
	Kind() string
	To() string
	Subject() string
}

type BaseMailer struct {
	kind string

	// Template variables
	ServiceName            string
	ServiceLogo            string
	ServiceUrl             string
	ServicePostalAddress   string
	ServiceCopyrightNotice string
}

func NewBaseMailer(kind string) *BaseMailer {
	return &BaseMailer{
		kind: kind,

		ServiceName:            viper.GetString("service_name"),
		ServiceLogo:            viper.GetString("service_logo"),
		ServiceUrl:             viper.GetString("service_url"),
		ServicePostalAddress:   viper.GetString("service_postal_address"),
		ServiceCopyrightNotice: viper.GetString("service_copyright_notice"),
	}
}

func (mailer *BaseMailer) Kind() string {
	return mailer.kind
}
