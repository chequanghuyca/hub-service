package repository

import (
	"context"
	"hub-service/module/email/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const emailCollection = "email_logs"

type EmailRepository interface {
	SaveEmailLog(ctx context.Context, email *model.EmailMessage) error
	GetEmailByID(ctx context.Context, id primitive.ObjectID) (*model.EmailMessage, error)
	UpdateEmailStatus(ctx context.Context, id primitive.ObjectID, status string, errorMsg string) error
	GetEmailLogs(ctx context.Context, page, limit int64) ([]model.EmailMessage, int64, error)
	GetEmailsByStatus(ctx context.Context, status string, limit int64) ([]model.EmailMessage, error)
	IncrementRetryCount(ctx context.Context, id primitive.ObjectID) error
}

type emailRepository struct {
	db *mongo.Database
}

func NewEmailRepository(db *mongo.Database) EmailRepository {
	return &emailRepository{db: db}
}

func (r *emailRepository) collection() *mongo.Collection {
	return r.db.Collection(emailCollection)
}

// SaveEmailLog saves a new email log to the database
func (r *emailRepository) SaveEmailLog(ctx context.Context, email *model.EmailMessage) error {
	if email.ID.IsZero() {
		email.ID = primitive.NewObjectID()
	}
	email.CreatedAt = time.Now()
	email.UpdatedAt = time.Now()

	_, err := r.collection().InsertOne(ctx, email)
	return err
}

// GetEmailByID retrieves an email by its ID
func (r *emailRepository) GetEmailByID(ctx context.Context, id primitive.ObjectID) (*model.EmailMessage, error) {
	var email model.EmailMessage
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &email, nil
}

// UpdateEmailStatus updates the status of an email
func (r *emailRepository) UpdateEmailStatus(ctx context.Context, id primitive.ObjectID, status string, errorMsg string) error {
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}

	if status == model.EmailStatusSent {
		now := time.Now()
		update["$set"].(bson.M)["sentAt"] = &now
	}

	if errorMsg != "" {
		update["$set"].(bson.M)["error"] = errorMsg
	}

	_, err := r.collection().UpdateByID(ctx, id, update)
	return err
}

// GetEmailLogs retrieves paginated email logs
func (r *emailRepository) GetEmailLogs(ctx context.Context, page, limit int64) ([]model.EmailMessage, int64, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var emails []model.EmailMessage
	if err := cursor.All(ctx, &emails); err != nil {
		return nil, 0, err
	}

	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return emails, total, nil
}

// GetEmailsByStatus retrieves emails by status
func (r *emailRepository) GetEmailsByStatus(ctx context.Context, status string, limit int64) ([]model.EmailMessage, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.M{"createdAt": 1})

	cursor, err := r.collection().Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var emails []model.EmailMessage
	if err := cursor.All(ctx, &emails); err != nil {
		return nil, err
	}

	return emails, nil
}

// IncrementRetryCount increments the retry count of an email
func (r *emailRepository) IncrementRetryCount(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$inc": bson.M{"retryCount": 1},
		"$set": bson.M{"updatedAt": time.Now()},
	}

	_, err := r.collection().UpdateByID(ctx, id, update)
	return err
}
