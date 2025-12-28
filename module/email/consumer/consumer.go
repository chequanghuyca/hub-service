package consumer

import (
	"context"
	"encoding/json"
	"hub-service/infrastructure/database/redis"
	"hub-service/infrastructure/messaging/kafka"
	"hub-service/module/email/model"
	"hub-service/module/email/repository"
	"hub-service/module/email/sender"
	"log"
	"os"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	defaultEmailTopic   = "email-notifications"
	defaultGroupID      = "email-service-group"
	emailCachePrefix    = "email:sent:"
	emailCacheTTL       = 24 * time.Hour
	maxRetries          = 3
	retryBackoffBase    = time.Second * 5
)

type EmailConsumer struct {
	kafkaClient *kafka.KafkaClient
	redisClient *redis.RedisClient
	repo        repository.EmailRepository
	sender      sender.EmailSender
	topic       string
	groupID     string
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewEmailConsumer(
	kafkaClient *kafka.KafkaClient,
	redisClient *redis.RedisClient,
	repo repository.EmailRepository,
	emailSender sender.EmailSender,
) *EmailConsumer {
	topic := os.Getenv("KAFKA_EMAIL_TOPIC")
	if topic == "" {
		topic = defaultEmailTopic
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = defaultGroupID
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EmailConsumer{
		kafkaClient: kafkaClient,
		redisClient: redisClient,
		repo:        repo,
		sender:      emailSender,
		topic:       topic,
		groupID:     groupID,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start starts the email consumer
func (c *EmailConsumer) Start() {
	if c.kafkaClient == nil {
		log.Println("Kafka client not available, email consumer not started")
		return
	}

	log.Printf("Starting email consumer for topic: %s, group: %s", c.topic, c.groupID)

	c.kafkaClient.StartConsumer(c.ctx, c.topic, c.groupID, c.handleMessage)
}

// Stop stops the email consumer
func (c *EmailConsumer) Stop() {
	log.Println("Stopping email consumer...")
	c.cancel()
}

// handleMessage processes a Kafka message
func (c *EmailConsumer) handleMessage(msg *kafkago.Message) error {
	var kafkaMsg model.KafkaEmailMessage
	if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
		log.Printf("Failed to unmarshal Kafka message: %v", err)
		return err
	}

	log.Printf("Processing email: %s to %v", kafkaMsg.EmailID, kafkaMsg.To)

	// Check if already sent (deduplication via Redis)
	if c.redisClient != nil {
		cacheKey := c.generateCacheKey(&kafkaMsg)
		exists, err := c.redisClient.Exists(cacheKey)
		if err != nil {
			log.Printf("Redis check error: %v", err)
		} else if exists {
			log.Printf("Email %s already sent, skipping", kafkaMsg.EmailID)
			return nil
		}
	}

	// Convert to EmailMessage for sending
	email := &model.EmailMessage{
		To:       kafkaMsg.To,
		Cc:       kafkaMsg.Cc,
		Bcc:      kafkaMsg.Bcc,
		From:     kafkaMsg.From,
		Subject:  kafkaMsg.Subject,
		Body:     kafkaMsg.Body,
		HTMLBody: kafkaMsg.HTMLBody,
		Priority: kafkaMsg.Priority,
	}

	// Send email with retry logic
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := retryBackoffBase * time.Duration(1<<uint(attempt-1))
			log.Printf("Retrying email %s (attempt %d/%d) after %v", kafkaMsg.EmailID, attempt+1, maxRetries, backoff)
			time.Sleep(backoff)
		}

		if err := c.sender.Send(email); err != nil {
			lastErr = err
			log.Printf("Failed to send email %s: %v", kafkaMsg.EmailID, err)

			// Update retry count in database
			if objID, err := primitive.ObjectIDFromHex(kafkaMsg.EmailID); err == nil {
				if err := c.repo.IncrementRetryCount(c.ctx, objID); err != nil {
					log.Printf("Failed to increment retry count: %v", err)
				}
			}
			continue
		}

		// Email sent successfully
		log.Printf("Email %s sent successfully", kafkaMsg.EmailID)

		// Update status in database
		if objID, err := primitive.ObjectIDFromHex(kafkaMsg.EmailID); err == nil {
			if err := c.repo.UpdateEmailStatus(c.ctx, objID, model.EmailStatusSent, ""); err != nil {
				log.Printf("Failed to update email status: %v", err)
			}
		}

		// Mark as sent in Redis for deduplication
		if c.redisClient != nil {
			cacheKey := c.generateCacheKey(&kafkaMsg)
			if _, err := c.redisClient.SetNX(cacheKey, "1", emailCacheTTL); err != nil {
				log.Printf("Failed to cache email: %v", err)
			}
		}

		return nil
	}

	// All retries failed
	log.Printf("Email %s failed after %d attempts", kafkaMsg.EmailID, maxRetries)
	if objID, err := primitive.ObjectIDFromHex(kafkaMsg.EmailID); err == nil {
		if err := c.repo.UpdateEmailStatus(c.ctx, objID, model.EmailStatusFailed, lastErr.Error()); err != nil {
			log.Printf("Failed to update email status: %v", err)
		}
	}

	return lastErr
}

// generateCacheKey generates a cache key for deduplication
func (c *EmailConsumer) generateCacheKey(email *model.KafkaEmailMessage) string {
	return emailCachePrefix + email.EmailID
}
