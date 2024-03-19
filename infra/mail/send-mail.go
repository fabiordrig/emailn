package mail

import (
	"emailn/domains/campaign"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(campaign *campaign.Campaign) error {
	m := gomail.NewMessage()

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 587, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASSWORD"))

	var emails []string

	for _, contact := range campaign.Contacts {
		emails = append(emails, contact.Email)
	}

	m.SetHeader("From", os.Getenv("MAIL_USER"))
	m.SetHeader("Subject", campaign.Name)
	m.SetHeader("To", emails...)

	m.SetBody("text/html", campaign.Content)

	return d.DialAndSend(m)

}
