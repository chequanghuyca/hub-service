package storagemail

import (
	emailmodel "hub-service/module/email/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func ResponseMeEmail(message string) error {
	godotenv.Load()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SYSTEM_EMAIL"))
	mailer.SetHeader("To", os.Getenv("SYSTEM_EMAIL"))
	mailer.SetHeader("Subject", "Có người xem Portfolio")
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		465,
		os.Getenv("SYSTEM_EMAIL"),
		os.Getenv("SYSTEM_EMAIL_SERVER"),
	)

	log.Println("Sending email to", dialer)

	if err := dialer.DialAndSend(mailer); err != nil {
		emailmodel.ErrSendEmail(err)
	}

	log.Println("Email sent successfully to", os.Getenv("SYSTEM_EMAIL"))
	return nil
}
