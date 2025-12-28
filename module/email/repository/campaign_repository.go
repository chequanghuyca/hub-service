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

const campaignCollection = "email_campaigns"

type CampaignRepository interface {
	Create(ctx context.Context, campaign *model.Campaign) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Campaign, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, page, limit int64) ([]*model.Campaign, int64, error)
	GetPendingCampaigns(ctx context.Context, beforeTime time.Time) ([]*model.Campaign, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string, err string) error
	UpdateEmailCounts(ctx context.Context, id primitive.ObjectID, total, sent, failed int) error
}

type campaignRepository struct {
	db *mongo.Database
}

func NewCampaignRepository(db *mongo.Database) CampaignRepository {
	return &campaignRepository{db: db}
}

func (r *campaignRepository) collection() *mongo.Collection {
	return r.db.Collection(campaignCollection)
}

// Create saves a new campaign to database
func (r *campaignRepository) Create(ctx context.Context, campaign *model.Campaign) error {
	if campaign.ID.IsZero() {
		campaign.ID = primitive.NewObjectID()
	}
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()
	campaign.Status = model.CampaignStatusPending

	_, err := r.collection().InsertOne(ctx, campaign)
	return err
}

// GetByID retrieves a campaign by ID
func (r *campaignRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Campaign, error) {
	var campaign model.Campaign
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&campaign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &campaign, nil
}

// Update updates a campaign with custom update document
func (r *campaignRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["$set"].(bson.M)["updated_at"] = time.Now()
	_, err := r.collection().UpdateByID(ctx, id, update)
	return err
}

// Delete removes a campaign
func (r *campaignRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List retrieves paginated campaigns
func (r *campaignRepository) List(ctx context.Context, page, limit int64) ([]*model.Campaign, int64, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var campaigns []*model.Campaign
	if err := cursor.All(ctx, &campaigns); err != nil {
		return nil, 0, err
	}

	total, err := r.collection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
}

// GetPendingCampaigns retrieves campaigns that are pending and scheduled before given time
func (r *campaignRepository) GetPendingCampaigns(ctx context.Context, beforeTime time.Time) ([]*model.Campaign, error) {
	filter := bson.M{
		"status":       model.CampaignStatusPending,
		"scheduled_at": bson.M{"$lte": beforeTime},
	}

	opts := options.Find().SetSort(bson.M{"scheduled_at": 1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var campaigns []*model.Campaign
	if err := cursor.All(ctx, &campaigns); err != nil {
		return nil, err
	}

	return campaigns, nil
}

// UpdateStatus updates campaign status
func (r *campaignRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string, errMsg string) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	if status == model.CampaignStatusCompleted || status == model.CampaignStatusFailed {
		now := time.Now()
		update["$set"].(bson.M)["processed_at"] = &now
	}

	if errMsg != "" {
		update["$set"].(bson.M)["error"] = errMsg
	}

	_, err := r.collection().UpdateByID(ctx, id, update)
	return err
}

// UpdateEmailCounts updates the email count statistics
func (r *campaignRepository) UpdateEmailCounts(ctx context.Context, id primitive.ObjectID, total, sent, failed int) error {
	update := bson.M{
		"$set": bson.M{
			"total_emails":  total,
			"sent_emails":   sent,
			"failed_emails": failed,
			"updated_at":    time.Now(),
		},
	}

	_, err := r.collection().UpdateByID(ctx, id, update)
	return err
}
