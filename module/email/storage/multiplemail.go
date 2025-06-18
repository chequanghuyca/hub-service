package storagemail

import (
	emailmodel "hub-service/module/email/model"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func MultipleSendEmail(listAddressesTo emailmodel.ListDataEmail, subject, body string) error {
	godotenv.Load()
	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		465,
		os.Getenv("SYSTEM_EMAIL"),
		os.Getenv("SYSTEM_EMAIL_SERVER"),
	)

	s, err := dialer.Dial()

	if err != nil {
		panic(err)
	}

	mailer := gomail.NewMessage()

	for _, address := range listAddressesTo {
		replacedBody := strings.Replace(body, "${name}", address.Name, -1)

		mailer.SetHeader("From", os.Getenv("SYSTEM_EMAIL"))
		mailer.SetHeader("To", address.Email)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", replacedBody)

		log.Println("Sending email to", replacedBody)

		if err := gomail.Send(s, mailer); err != nil {
			log.Printf("Could not send email to %q: %v", address.Email, err)
		}

		log.Println("Email sent successfully to", address.Email)

		mailer.Reset()
	}

	return nil
}
