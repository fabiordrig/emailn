package mail

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail() {
	m := gomail.NewMessage()

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), 587, os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASSWORD"))

	m.SetHeader("From", os.Getenv("MAIL_USER"))
	m.SetHeader("To", "fabiordrig@gmail.com")

	m.SetBody("text/html", "Hello <b>World</b>")

	err := d.DialAndSend(m)

	if err != nil {
		println("Error sending email", err.Error())
		panic(err)
	}

}
