package mailers

import (
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/aymerick/douceur/inliner"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

type Sender struct {
	mailer   Mailer
	smtpConf *SMTPConf
	noop     bool
}

type SMTPConf struct {
	From string
	Host string
	Port int
	User string
	Pass string
}

func NewSender(mailer Mailer) *Sender {
	return &Sender{
		mailer: mailer,
		smtpConf: &SMTPConf{
			From: viper.GetString("smtp_from"),
			Host: viper.GetString("smtp_host"),
			Port: viper.GetInt("smtp_port"),
			User: viper.GetString("smtp_auth_user"),
			Pass: viper.GetString("smtp_auth_pass"),
		},
	}
}

func (sender *Sender) SetSMTPConf(conf *SMTPConf) {
	sender.smtpConf = conf
}

func (sender *Sender) SetNoop(isNoop bool) {
	sender.noop = isNoop
}

func (sender *Sender) Send() error {
	var result error

	mail := sender.newEmail()

	if !sender.noop {
		result = mail.Send(sender.smtpAddr(), sender.smtpAuth())
	}

	return result
}

func (sender *Sender) newEmail() *email.Email {
	// generate HTML
	rawHtml := sender.content(TPL_HTML)

	// inline CSS
	htmlContent, err := inliner.Inline(rawHtml)
	if err != nil {
		panic(err)
	}

	return &email.Email{
		To:      []string{sender.mailer.To()},
		From:    sender.smtpConf.From,
		Subject: sender.mailer.Subject(),
		HTML:    []byte(htmlContent),
		Text:    []byte(sender.content(TPL_TEXT)),
		Headers: textproto.MIMEHeader{},
	}
}

func (sender *Sender) content(tplKind TplKind) string {
	result, err := templater.Generate(sender.mailer.Kind(), tplKind, sender.mailer)
	if err != nil {
		panic(err)
	}

	return result
}

func (sender *Sender) smtpAddr() string {
	return fmt.Sprintf("%s:%d", sender.smtpConf.Host, sender.smtpConf.Port)
}

func (sender *Sender) smtpAuth() smtp.Auth {
	if sender.smtpConf.User != "" {
		return smtp.PlainAuth("", sender.smtpConf.User, sender.smtpConf.Pass, sender.smtpConf.Host)
	} else {
		return nil
	}
}
