package notifications

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

/*
	* settings for email smtp server
	例如：smtp.qq.com:465
	SmtpServer string  smtp.qq.com
	SmtpPort   int     465
	SmtpUsername string  sender email address xxxx@qq.com
	smtpPassword string  sender email password xxxxxxxx
*/
type EmailxOptions struct {
	SmtpServer   string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	TLSConfig    *tls.Config
}

type Emailx struct {
	dialer  *gomail.Dialer
	options EmailxOptions
}

func NewEmail(options EmailxOptions) *Emailx {
	return &Emailx{
		dialer:  gomail.NewDialer(options.SmtpServer, options.SmtpPort, options.SmtpUsername, options.SmtpPassword),
		options: options}
}

func (e *Emailx) Send(to string, subject string, buf []byte) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.dialer.Username)
	m.SetAddressHeader("To", to, to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", string(buf))

	_ = m
	err := e.dialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
