package service

import (
	storagemail "hub-service/module/email/storage"
	"hub-service/module/email/template"
	"log"
	"os"
)

type WelcomeEmailService struct{}

func NewWelcomeEmailService() *WelcomeEmailService {
	return &WelcomeEmailService{}
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *WelcomeEmailService) SendWelcomeEmail(userName, userEmail string) error {
	// Get login URL from environment variable or use default
	loginUrl := os.Getenv("BASE_URL_TRANSMASTER_PROD")
	if loginUrl == "" {
		loginUrl = "https://your-app.com/login" // Default fallback
	}

	// Prepare email data
	emailData := template.MailWelcomeData{
		Name:     userName,
		LoginUrl: loginUrl,
	}

	// Generate email content
	subject := template.GetSubjectMailWelcome()
	body := template.GetBodyMailWelcome(emailData)

	// Send email
	err := storagemail.SingleSendEmail(userEmail, subject, body)
	if err != nil {
		log.Printf("Failed to send welcome email to %s: %v", userEmail, err)
		return err
	}

	log.Printf("Welcome email sent successfully to %s", userEmail)
	return nil
}
