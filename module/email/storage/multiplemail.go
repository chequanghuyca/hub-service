package storagemail

import (
	emailmodel "hub-service/module/email/model"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func MultipleSendEmail(listAddressesTo emailmodel.ListDataEmail, subject, body string) error {
	godotenv.Load()

	// Get environment variables
	smtpHost := os.Getenv("SYSTEM_EMAIL_HOST")
	smtpPortStr := os.Getenv("SYSTEM_EMAIL_PORT")
	smtpEmail := os.Getenv("SYSTEM_EMAIL")
	smtpPassword := os.Getenv("SYSTEM_EMAIL_SERVER")

	// Convert port to integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("Invalid SYSTEM_EMAIL_PORT: %s, using default 587", smtpPortStr)
		smtpPort = 587 // Use 587 as default for better DigitalOcean compatibility
	}

	// Default values if not set
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}

	log.Printf("Using SMTP config for multiple emails - Host: %s, Port: %d, Email: %s", smtpHost, smtpPort, smtpEmail)

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpEmail,
		smtpPassword,
	)

	s, err := dialer.Dial()

	if err != nil {
		panic(err)
	}

	mailer := gomail.NewMessage()

	for _, address := range listAddressesTo {
		replacedBody := strings.Replace(body, "${name}", address.Name, -1)

		mailer.SetHeader("From", smtpEmail)
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
