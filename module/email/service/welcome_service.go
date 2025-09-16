package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type WelcomeEmailService struct{}

func NewWelcomeEmailService() *WelcomeEmailService {
	return &WelcomeEmailService{}
}

// SendWelcomeEmail sends a welcome email to a new user via external email service
func (s *WelcomeEmailService) SendWelcomeEmail(userName, userEmail string) error {
	godotenv.Load()

	// Get external email service configuration
	baseUrl := os.Getenv("SYSTEM_EMAIL_SERVICE_BASE_URL")
	apiKey := os.Getenv("SYSTEM_EMAIL_SERVICE_API_KEY")
	loginUrl := os.Getenv("BASE_URL_TRANSMASTER_PROD")

	if baseUrl == "" {
		return fmt.Errorf("SYSTEM_EMAIL_SERVICE_BASE_URL environment variable is required")
	}
	if apiKey == "" {
		return fmt.Errorf("SYSTEM_EMAIL_SERVICE_API_KEY environment variable is required")
	}
	if loginUrl == "" {
		loginUrl = "https://transmaster.site"
	}

	// Prepare request payload
	requestData := map[string]string{
		"email":    userEmail,
		"name":     userName,
		"loginUrl": loginUrl,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Create HTTP request
	url := baseUrl + "/api/email/welcome-user"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to call external email service: %v", err)
		return fmt.Errorf("failed to call external email service: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("External email service returned error: status %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("external email service error: status %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Welcome email sent successfully to %s via external service", userEmail)
	return nil
}
