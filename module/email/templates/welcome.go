package templates

import (
	"log"
)

// WelcomeEmailData contains data for welcome email template
type WelcomeEmailData struct {
	Name     string
	LoginUrl string
}

// GetWelcomeEmailSubject returns the subject for welcome email
func GetWelcomeEmailSubject() string {
	return "Welcome to TransMaster!"
}

// GetWelcomeEmailHTML returns the HTML body for welcome email
// It uses the external HTML template file for easier customization
func GetWelcomeEmailHTML(data WelcomeEmailData) string {
	html, err := RenderTemplate("welcome", data)
	if err != nil {
		log.Printf("Error: Failed to render welcome template: %v", err)
		return ""
	}
	return html
}
