package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	EmailStatusPending   = "pending"
	EmailStatusQueued    = "queued"
	EmailStatusSent      = "sent"
	EmailStatusFailed    = "failed"
	EmailStatusDelivered = "delivered"

	EmailPriorityHigh   = "high"
	EmailPriorityNormal = "normal"
	EmailPriorityLow    = "low"
)

// EmailMessage represents an email to be sent
type EmailMessage struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	To          []string               `json:"to" bson:"to"`
	Cc          []string               `json:"cc,omitempty" bson:"cc,omitempty"`
	Bcc         []string               `json:"bcc,omitempty" bson:"bcc,omitempty"`
	From        string                 `json:"from,omitempty" bson:"from,omitempty"`
	Subject     string                 `json:"subject" bson:"subject"`
	Body        string                 `json:"body,omitempty" bson:"body,omitempty"`
	HTMLBody    string                 `json:"htmlBody,omitempty" bson:"htmlBody,omitempty"`
	Template    string                 `json:"template,omitempty" bson:"template,omitempty"`
	TemplateData map[string]interface{} `json:"templateData,omitempty" bson:"templateData,omitempty"`
	Priority    string                 `json:"priority" bson:"priority"`
	Status      string                 `json:"status" bson:"status"`
	SentAt      *time.Time             `json:"sentAt,omitempty" bson:"sentAt,omitempty"`
	Error       string                 `json:"error,omitempty" bson:"error,omitempty"`
	RetryCount  int                    `json:"retryCount" bson:"retryCount"`
	MaxRetries  int                    `json:"maxRetries" bson:"maxRetries"`
	CreatedAt   time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updatedAt"`
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Subject     string             `json:"subject" bson:"subject"`
	HTMLContent string             `json:"htmlContent" bson:"htmlContent"`
	TextContent string             `json:"textContent,omitempty" bson:"textContent,omitempty"`
	Variables   []string           `json:"variables" bson:"variables"`
	IsActive    bool               `json:"isActive" bson:"isActive"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// SendEmailRequest represents an API request to send email
type SendEmailRequest struct {
	To           []string               `json:"to" binding:"required"`
	Cc           []string               `json:"cc,omitempty"`
	Bcc          []string               `json:"bcc,omitempty"`
	Subject      string                 `json:"subject" binding:"required"`
	Body         string                 `json:"body,omitempty"`
	HTMLBody     string                 `json:"htmlBody,omitempty"`
	Template     string                 `json:"template,omitempty"`
	TemplateData map[string]interface{} `json:"templateData,omitempty"`
	Priority     string                 `json:"priority,omitempty"`
}

// SendBulkEmailRequest represents a request to send bulk emails
type SendBulkEmailRequest struct {
	Emails []SendEmailRequest `json:"emails" binding:"required"`
}

// EmailResponse represents the API response
type EmailResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// KafkaEmailMessage is the message format for Kafka
type KafkaEmailMessage struct {
	EmailID   string                 `json:"emailId"`
	To        []string               `json:"to"`
	Cc        []string               `json:"cc,omitempty"`
	Bcc       []string               `json:"bcc,omitempty"`
	From      string                 `json:"from,omitempty"`
	Subject   string                 `json:"subject"`
	Body      string                 `json:"body,omitempty"`
	HTMLBody  string                 `json:"htmlBody,omitempty"`
	Priority  string                 `json:"priority"`
	Timestamp time.Time              `json:"timestamp"`
}
