package storagemail

import (
	emailmodel "hub-service/module/email/model"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func ResponseMeEmail(message string) error {
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

	log.Printf("Using SMTP config for notification - Host: %s, Port: %d, Email: %s", smtpHost, smtpPort, smtpEmail)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpEmail)
	mailer.SetHeader("To", smtpEmail)
	mailer.SetHeader("Subject", "Có người xem Portfolio")
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpEmail,
		smtpPassword,
	)

	log.Println("Sending email to", dialer)

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send notification email: %v", err)
		return emailmodel.ErrSendEmail(err)
	}

	log.Println("Email sent successfully to", os.Getenv("SYSTEM_EMAIL"))
	return nil
}
