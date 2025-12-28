package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Campaign status constants
const (
	CampaignStatusPending    = "pending"
	CampaignStatusProcessing = "processing"
	CampaignStatusCompleted  = "completed"
	CampaignStatusCancelled  = "cancelled"
	CampaignStatusFailed     = "failed"
)

// Campaign represents a scheduled email campaign
type Campaign struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Subject      string             `json:"subject" bson:"subject"`
	HTMLBody     string             `json:"html_body" bson:"html_body"`
	ScheduledAt  time.Time          `json:"scheduled_at" bson:"scheduled_at"`
	Status       string             `json:"status" bson:"status"`
	TestEmails   []string           `json:"test_emails,omitempty" bson:"test_emails,omitempty"`
	TotalEmails  int                `json:"total_emails" bson:"total_emails"`
	SentEmails   int                `json:"sent_emails" bson:"sent_emails"`
	FailedEmails int                `json:"failed_emails" bson:"failed_emails"`
	CreatedBy    primitive.ObjectID `json:"created_by" bson:"created_by"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	ProcessedAt  *time.Time         `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	Error        string             `json:"error,omitempty" bson:"error,omitempty"`
}

// IsTestMode returns true if campaign is in test mode (has test_emails)
func (c *Campaign) IsTestMode() bool {
	return len(c.TestEmails) > 0
}

// CreateCampaignRequest represents API request to create a campaign
type CreateCampaignRequest struct {
	Subject     string   `json:"subject" binding:"required"`
	HTMLBody    string   `json:"html_body" binding:"required"`
	ScheduledAt int64    `json:"scheduled_at" binding:"required"` // Unix timestamp in seconds
	TestEmails  []string `json:"test_emails,omitempty"`
}

// GetScheduledTime converts Unix timestamp to time.Time
func (r *CreateCampaignRequest) GetScheduledTime() time.Time {
	return time.Unix(r.ScheduledAt, 0)
}

// UpdateCampaignRequest represents API request to update a campaign
type UpdateCampaignRequest struct {
	Subject     *string  `json:"subject,omitempty"`
	HTMLBody    *string  `json:"html_body,omitempty"`
	ScheduledAt *int64   `json:"scheduled_at,omitempty"` // Unix timestamp in seconds
	TestEmails  []string `json:"test_emails,omitempty"`
}

// GetScheduledTime converts Unix timestamp to time.Time, returns nil if not set
func (r *UpdateCampaignRequest) GetScheduledTime() *time.Time {
	if r.ScheduledAt == nil {
		return nil
	}
	t := time.Unix(*r.ScheduledAt, 0)
	return &t
}

// CampaignResponse represents API response for a campaign
type CampaignResponse struct {
	ID           string     `json:"id"`
	Subject      string     `json:"subject"`
	HTMLBody     string     `json:"html_body"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	Status       string     `json:"status"`
	TestMode     bool       `json:"test_mode"`
	TestEmails   []string   `json:"test_emails,omitempty"`
	TotalEmails  int        `json:"total_emails"`
	SentEmails   int        `json:"sent_emails"`
	FailedEmails int        `json:"failed_emails"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	Error        string     `json:"error,omitempty"`
}

// ToCampaignResponse converts Campaign to CampaignResponse
func (c *Campaign) ToCampaignResponse() *CampaignResponse {
	return &CampaignResponse{
		ID:           c.ID.Hex(),
		Subject:      c.Subject,
		HTMLBody:     c.HTMLBody,
		ScheduledAt:  c.ScheduledAt,
		Status:       c.Status,
		TestMode:     c.IsTestMode(),
		TestEmails:   c.TestEmails,
		TotalEmails:  c.TotalEmails,
		SentEmails:   c.SentEmails,
		FailedEmails: c.FailedEmails,
		CreatedBy:    c.CreatedBy.Hex(),
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		ProcessedAt:  c.ProcessedAt,
		Error:        c.Error,
	}
}

// CampaignListResponse represents paginated campaign list response
type CampaignListResponse struct {
	Campaigns  []*CampaignResponse `json:"campaigns"`
	Total      int64               `json:"total"`
	Page       int64               `json:"page"`
	Limit      int64               `json:"limit"`
	TotalPages int64               `json:"total_pages"`
}
