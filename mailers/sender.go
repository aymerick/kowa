package mailers

import (
	"fmt"
	"log"
	"net/smtp"
	"net/textproto"

	"github.com/aymerick/douceur/inliner"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

// Sender send mail via SMTP
type Sender struct {
	mailer   Mailer
	smtpConf *SMTPConf
	noop     bool
}

// SMTPConf holds SMTP configuration
type SMTPConf struct {
	From string
	Host string
	Port int
	User string
	Pass string
}

// NewSender instanciates a new Sender
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

// SetSMTPConf sets the SMTP configuration
func (sender *Sender) SetSMTPConf(conf *SMTPConf) {
	sender.smtpConf = conf
}

// SetNoop sets sender in NOOP mode
func (sender *Sender) SetNoop(isNoop bool) {
	sender.noop = isNoop
}

// Send triggers mail sending
func (sender *Sender) Send() error {
	var result error

	mail := sender.newEmail()

	if !sender.noop {
		log.Printf("Sending email to: %v", mail.To)
		result = mail.Send(sender.smtpAddr(), sender.smtpAuth())
		if result == nil {
			log.Printf("Mail successfully sent to: %v", mail.To)
		} else {
			log.Printf("Failed to send email to: %v - %s", mail.To, result)
		}
	}

	return result
}

func (sender *Sender) newEmail() *email.Email {
	// generate HTML
	rawHTML := sender.content(tplHTML)

	// inline CSS
	htmlContent, err := inliner.Inline(rawHTML)
	if err != nil {
		panic(err)
	}

	return &email.Email{
		To:      []string{sender.mailer.To()},
		From:    sender.smtpConf.From,
		Subject: sender.mailer.Subject(),
		HTML:    []byte(htmlContent),
		Text:    []byte(sender.content(tplText)),
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
	}

	return nil
}
