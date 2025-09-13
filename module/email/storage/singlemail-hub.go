package storagemail

import (
	emailmodel "hub-service/module/email/model"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SingleSendEmailHub(to string, subject string, body string) error {
	godotenv.Load()

	// Get environment variables for TransMaster
	smtpHost := os.Getenv("SYSTEM_EMAIL_HOST")
	smtpPortStr := os.Getenv("SYSTEM_EMAIL_PORT")
	smtpEmail := os.Getenv("SYSTEM_EMAIL_TRANSMASTER")
	smtpPassword := os.Getenv("SYSTEM_EMAIL_SERVER_TRANSMASTER")

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

	log.Printf("Using SMTP config for TransMaster - Host: %s, Port: %d, Email: %s", smtpHost, smtpPort, smtpEmail)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "TransMaster <"+smtpEmail+">")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpEmail,
		smtpPassword,
	)

	log.Println("Sending email to", dialer)

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send TransMaster email to %s: %v", to, err)
		return emailmodel.ErrSendEmail(err)
	}

	log.Println("TransMaster email sent successfully to", to)
	return nil
}
