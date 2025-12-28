package biz

import (
	"context"
	"fmt"
	"hub-service/module/email/model"
	"hub-service/module/email/repository"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CampaignBusiness interface {
	CreateCampaign(ctx context.Context, req *model.CreateCampaignRequest, createdBy primitive.ObjectID) (*model.CampaignResponse, error)
	GetCampaign(ctx context.Context, id string) (*model.CampaignResponse, error)
	ListCampaigns(ctx context.Context, page, limit int64) (*model.CampaignListResponse, error)
	UpdateCampaign(ctx context.Context, id string, req *model.UpdateCampaignRequest) (*model.CampaignResponse, error)
	CancelCampaign(ctx context.Context, id string) error
	GetPendingCampaigns(ctx context.Context) ([]*model.Campaign, error)
	ProcessCampaign(ctx context.Context, campaign *model.Campaign, userEmails []string) error
}

type campaignBusiness struct {
	repo         repository.CampaignRepository
	emailBiz     EmailBusiness
}

func NewCampaignBusiness(repo repository.CampaignRepository, emailBiz EmailBusiness) CampaignBusiness {
	return &campaignBusiness{
		repo:     repo,
		emailBiz: emailBiz,
	}
}

// CreateCampaign creates a new scheduled campaign
func (b *campaignBusiness) CreateCampaign(ctx context.Context, req *model.CreateCampaignRequest, createdBy primitive.ObjectID) (*model.CampaignResponse, error) {
	scheduledTime := req.GetScheduledTime()
	
	// Validate scheduled time is in the future
	if scheduledTime.Before(time.Now()) {
		return nil, fmt.Errorf("scheduled_at must be in the future")
	}

	campaign := &model.Campaign{
		Subject:     req.Subject,
		HTMLBody:    req.HTMLBody,
		ScheduledAt: scheduledTime,
		TestEmails:  req.TestEmails,
		CreatedBy:   createdBy,
		Status:      model.CampaignStatusPending,
	}

	if err := b.repo.Create(ctx, campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	log.Printf("Campaign created: %s, scheduled for: %s, test_mode: %v", 
		campaign.ID.Hex(), campaign.ScheduledAt.Format(time.RFC3339), campaign.IsTestMode())

	return campaign.ToCampaignResponse(), nil
}

// GetCampaign retrieves a campaign by ID
func (b *campaignBusiness) GetCampaign(ctx context.Context, id string) (*model.CampaignResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID: %w", err)
	}

	campaign, err := b.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, nil
	}

	return campaign.ToCampaignResponse(), nil
}

// ListCampaigns retrieves paginated campaigns
func (b *campaignBusiness) ListCampaigns(ctx context.Context, page, limit int64) (*model.CampaignListResponse, error) {
	campaigns, total, err := b.repo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.CampaignResponse, len(campaigns))
	for i, c := range campaigns {
		responses[i] = c.ToCampaignResponse()
	}

	totalPages := (total + limit - 1) / limit

	return &model.CampaignListResponse{
		Campaigns:  responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateCampaign updates a pending campaign
func (b *campaignBusiness) UpdateCampaign(ctx context.Context, id string, req *model.UpdateCampaignRequest) (*model.CampaignResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID: %w", err)
	}

	// Check campaign exists and is pending
	campaign, err := b.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, fmt.Errorf("campaign not found")
	}
	if campaign.Status != model.CampaignStatusPending {
		return nil, fmt.Errorf("can only update pending campaigns")
	}

	// Build update document
	updateData := bson.M{}
	if req.Subject != nil {
		updateData["subject"] = *req.Subject
	}
	if req.HTMLBody != nil {
		updateData["html_body"] = *req.HTMLBody
	}
	if req.ScheduledAt != nil {
		scheduledTime := req.GetScheduledTime()
		if scheduledTime.Before(time.Now()) {
			return nil, fmt.Errorf("scheduled_at must be in the future")
		}
		updateData["scheduled_at"] = *scheduledTime
	}
	if req.TestEmails != nil {
		updateData["test_emails"] = req.TestEmails
	}

	if len(updateData) == 0 {
		return campaign.ToCampaignResponse(), nil
	}

	update := bson.M{"$set": updateData}
	if err := b.repo.Update(ctx, objectID, update); err != nil {
		return nil, err
	}

	// Get updated campaign
	campaign, _ = b.repo.GetByID(ctx, objectID)
	return campaign.ToCampaignResponse(), nil
}

// CancelCampaign cancels a pending campaign
func (b *campaignBusiness) CancelCampaign(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid campaign ID: %w", err)
	}

	campaign, err := b.repo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}
	if campaign == nil {
		return fmt.Errorf("campaign not found")
	}
	if campaign.Status != model.CampaignStatusPending {
		return fmt.Errorf("can only cancel pending campaigns")
	}

	if err := b.repo.UpdateStatus(ctx, objectID, model.CampaignStatusCancelled, ""); err != nil {
		return err
	}

	log.Printf("Campaign cancelled: %s", id)
	return nil
}

// GetPendingCampaigns returns campaigns ready to be processed
func (b *campaignBusiness) GetPendingCampaigns(ctx context.Context) ([]*model.Campaign, error) {
	return b.repo.GetPendingCampaigns(ctx, time.Now())
}

// ProcessCampaign processes a campaign and queues emails
func (b *campaignBusiness) ProcessCampaign(ctx context.Context, campaign *model.Campaign, userEmails []string) error {
	log.Printf("Processing campaign: %s", campaign.ID.Hex())

	// Update status to processing
	if err := b.repo.UpdateStatus(ctx, campaign.ID, model.CampaignStatusProcessing, ""); err != nil {
		return err
	}

	// Determine recipients
	var recipients []string
	if campaign.IsTestMode() {
		recipients = campaign.TestEmails
		log.Printf("Test mode: sending to %d test emails", len(recipients))
	} else {
		recipients = userEmails
		log.Printf("Production mode: sending to %d users", len(recipients))
	}

	// Update total count
	if err := b.repo.UpdateEmailCounts(ctx, campaign.ID, len(recipients), 0, 0); err != nil {
		log.Printf("Failed to update email counts: %v", err)
	}

	// Queue emails for each recipient
	sentCount := 0
	failedCount := 0

	for _, email := range recipients {
		req := &model.SendEmailRequest{
			To:       []string{email},
			Subject:  campaign.Subject,
			HTMLBody: campaign.HTMLBody,
			Priority: model.EmailPriorityNormal,
		}

		_, err := b.emailBiz.QueueEmail(ctx, req)
		if err != nil {
			log.Printf("Failed to queue email for %s: %v", email, err)
			failedCount++
		} else {
			sentCount++
		}
	}

	// Update final counts and status
	if err := b.repo.UpdateEmailCounts(ctx, campaign.ID, len(recipients), sentCount, failedCount); err != nil {
		log.Printf("Failed to update final email counts: %v", err)
	}

	finalStatus := model.CampaignStatusCompleted
	var errMsg string
	if failedCount > 0 && sentCount == 0 {
		finalStatus = model.CampaignStatusFailed
		errMsg = "All emails failed to queue"
	} else if failedCount > 0 {
		errMsg = fmt.Sprintf("%d emails failed to queue", failedCount)
	}

	if err := b.repo.UpdateStatus(ctx, campaign.ID, finalStatus, errMsg); err != nil {
		return err
	}

	log.Printf("Campaign %s completed: %d sent, %d failed", campaign.ID.Hex(), sentCount, failedCount)
	return nil
}
