package mailers

import (
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

type Sender struct {
	mailer   Mailer
	smtpConf *SMTPConf
}

// provides mail data to Sender
type Mailer interface {
	To() string
	Subject() string
	Html() string
	Text() string
}

type SMTPConf struct {
	from string
	host string
	port int
	user string
	pass string
}

func NewSender(mailer Mailer) *Sender {
	return &Sender{
		mailer: mailer,
		smtpConf: &SMTPConf{
			from: viper.GetString("smtp_from"),
			host: viper.GetString("smtp_host"),
			port: viper.GetInt("smtp_port"),
			user: viper.GetString("smtp_auth_user"),
			pass: viper.GetString("smtp_auth_pass"),
		},
	}
}

func (sender *Sender) smtpAddr() string {
	return fmt.Sprintf("%s:%d", sender.smtpConf.host, sender.smtpConf.port)
}

func (sender *Sender) auth() smtp.Auth {
	if sender.smtpConf.user != "" {
		return smtp.PlainAuth("", sender.smtpConf.user, sender.smtpConf.pass, sender.smtpConf.host)
	} else {
		return nil
	}
}

func (sender *Sender) NewEmail() *email.Email {
	return &email.Email{
		To:      []string{sender.mailer.To()},
		From:    sender.smtpConf.from,
		Subject: sender.mailer.Subject(),
		HTML:    []byte(sender.mailer.Html()),
		Text:    []byte(sender.mailer.Text()),
		Headers: textproto.MIMEHeader{},
	}
}

func (sender *Sender) Send() error {
	mail := sender.NewEmail()

	return mail.Send(sender.smtpAddr(), sender.auth())
}
