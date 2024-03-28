// mail/send.go
package mail

import (
	"emailn/domains/campaign"
	"os"

	"gopkg.in/gomail.v2"
)

type SMTPSender struct {
	Dialer *gomail.Dialer
}

func NewSMTPSender() *SMTPSender {
	return &SMTPSender{
		Dialer: gomail.NewDialer(os.Getenv("MAIL_HOST"), 587, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASSWORD")),
	}
}

func (s *SMTPSender) SendEmail(c *campaign.Campaign) error {
	m := gomail.NewMessage()

	var emails []string
	for _, contact := range c.Contacts {
		emails = append(emails, contact.Email)
	}

	m.SetHeader("From", os.Getenv("MAIL_USER"))
	m.SetHeader("Subject", c.Name)
	m.SetHeader("To", emails...)
	m.SetBody("text/html", c.Content)

	return s.Dialer.DialAndSend(m)
}
