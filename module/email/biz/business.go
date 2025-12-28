package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"hub-service/infrastructure/database/redis"
	"hub-service/infrastructure/messaging/kafka"
	"hub-service/module/email/model"
	"hub-service/module/email/repository"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	defaultEmailTopic = "email-notifications"
	emailCachePrefix  = "email:sent:"
	emailCacheTTL     = 24 * time.Hour
)

type EmailBusiness interface {
	QueueEmail(ctx context.Context, req *model.SendEmailRequest) (*model.EmailResponse, error)
	QueueBulkEmails(ctx context.Context, req *model.SendBulkEmailRequest) ([]model.EmailResponse, error)
	GetEmailLogs(ctx context.Context, page, limit int64) ([]model.EmailMessage, int64, error)
	GetEmailByID(ctx context.Context, id string) (*model.EmailMessage, error)
}

type emailBusiness struct {
	repo        repository.EmailRepository
	kafkaClient *kafka.KafkaClient
	redisClient *redis.RedisClient
	topic       string
}

func NewEmailBusiness(
	repo repository.EmailRepository,
	kafkaClient *kafka.KafkaClient,
	redisClient *redis.RedisClient,
) EmailBusiness {
	topic := os.Getenv("KAFKA_EMAIL_TOPIC")
	if topic == "" {
		topic = defaultEmailTopic
	}

	return &emailBusiness{
		repo:        repo,
		kafkaClient: kafkaClient,
		redisClient: redisClient,
		topic:       topic,
	}
}

// QueueEmail queues a single email for sending via Kafka
func (b *emailBusiness) QueueEmail(ctx context.Context, req *model.SendEmailRequest) (*model.EmailResponse, error) {
	// Create email message
	email := &model.EmailMessage{
		ID:           primitive.NewObjectID(),
		To:           req.To,
		Cc:           req.Cc,
		Bcc:          req.Bcc,
		Subject:      req.Subject,
		Body:         req.Body,
		HTMLBody:     req.HTMLBody,
		Template:     req.Template,
		TemplateData: req.TemplateData,
		Priority:     req.Priority,
		Status:       model.EmailStatusPending,
		MaxRetries:   3,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if email.Priority == "" {
		email.Priority = model.EmailPriorityNormal
	}

	// Save to database first
	if err := b.repo.SaveEmailLog(ctx, email); err != nil {
		return nil, fmt.Errorf("failed to save email log: %w", err)
	}

	// Check for duplicate (using Redis if available)
	if b.redisClient != nil {
		cacheKey := b.generateCacheKey(email)
		exists, err := b.redisClient.Exists(cacheKey)
		if err != nil {
			log.Printf("Redis check error: %v", err)
		} else if exists {
			return &model.EmailResponse{
				ID:      email.ID.Hex(),
				Status:  "duplicate",
				Message: "Email already sent recently",
			}, nil
		}
	}

	// Queue to Kafka
	if b.kafkaClient != nil {
		kafkaMsg := model.KafkaEmailMessage{
			EmailID:   email.ID.Hex(),
			To:        email.To,
			Cc:        email.Cc,
			Bcc:       email.Bcc,
			Subject:   email.Subject,
			Body:      email.Body,
			HTMLBody:  email.HTMLBody,
			Priority:  email.Priority,
			Timestamp: time.Now(),
		}

		msgBytes, err := json.Marshal(kafkaMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Kafka message: %w", err)
		}

		if err := b.kafkaClient.Produce(ctx, b.topic, []byte(email.ID.Hex()), msgBytes); err != nil {
			return nil, fmt.Errorf("failed to queue email to Kafka: %w", err)
		}

		// Update status to queued
		if err := b.repo.UpdateEmailStatus(ctx, email.ID, model.EmailStatusQueued, ""); err != nil {
			log.Printf("Failed to update email status: %v", err)
		}
	}

	return &model.EmailResponse{
		ID:      email.ID.Hex(),
		Status:  model.EmailStatusQueued,
		Message: "Email queued successfully",
	}, nil
}

// QueueBulkEmails queues multiple emails for sending
func (b *emailBusiness) QueueBulkEmails(ctx context.Context, req *model.SendBulkEmailRequest) ([]model.EmailResponse, error) {
	responses := make([]model.EmailResponse, 0, len(req.Emails))

	for _, emailReq := range req.Emails {
		resp, err := b.QueueEmail(ctx, &emailReq)
		if err != nil {
			responses = append(responses, model.EmailResponse{
				Status:  model.EmailStatusFailed,
				Message: err.Error(),
			})
			continue
		}
		responses = append(responses, *resp)
	}

	return responses, nil
}

// GetEmailLogs retrieves paginated email logs
func (b *emailBusiness) GetEmailLogs(ctx context.Context, page, limit int64) ([]model.EmailMessage, int64, error) {
	return b.repo.GetEmailLogs(ctx, page, limit)
}

// GetEmailByID retrieves an email by its ID
func (b *emailBusiness) GetEmailByID(ctx context.Context, id string) (*model.EmailMessage, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid email ID: %w", err)
	}
	return b.repo.GetEmailByID(ctx, objectID)
}

// generateCacheKey generates a cache key for deduplication
func (b *emailBusiness) generateCacheKey(email *model.EmailMessage) string {
	// Create a unique key based on recipients and subject
	return fmt.Sprintf("%s%s:%s", emailCachePrefix, email.To[0], email.Subject)
}
