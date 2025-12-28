package scheduler

import (
	"context"
	"hub-service/core/appctx"
	"hub-service/module/email/biz"
	"hub-service/module/email/repository"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	schedulerInterval = 1 * time.Minute
)

type CampaignScheduler struct {
	campaignBiz biz.CampaignBusiness
	appCtx      appctx.AppContext
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	running     bool
	mu          sync.Mutex
}

func NewCampaignScheduler(appCtx appctx.AppContext) *CampaignScheduler {
	db := appCtx.GetDatabase().MongoDB.Database
	
	// Initialize repositories
	emailRepo := repository.NewEmailRepository(db)
	campaignRepo := repository.NewCampaignRepository(db)
	
	// Initialize business layers
	emailBiz := biz.NewEmailBusiness(emailRepo, appCtx.GetKafka(), appCtx.GetRedis())
	campaignBiz := biz.NewCampaignBusiness(campaignRepo, emailBiz)
	
	ctx, cancel := context.WithCancel(context.Background())

	return &CampaignScheduler{
		campaignBiz: campaignBiz,
		appCtx:      appCtx,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start starts the scheduler background job
func (s *CampaignScheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	log.Println("Starting campaign scheduler...")

	s.wg.Add(1)
	go s.run()
}

// Stop stops the scheduler
func (s *CampaignScheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	log.Println("Stopping campaign scheduler...")
	s.cancel()
	s.wg.Wait()
	log.Println("Campaign scheduler stopped")
}

func (s *CampaignScheduler) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(schedulerInterval)
	defer ticker.Stop()

	// Run immediately on start
	s.processPendingCampaigns()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.processPendingCampaigns()
		}
	}
}

func (s *CampaignScheduler) processPendingCampaigns() {
	ctx := context.Background()

	campaigns, err := s.campaignBiz.GetPendingCampaigns(ctx)
	if err != nil {
		log.Printf("Error getting pending campaigns: %v", err)
		return
	}

	if len(campaigns) == 0 {
		return
	}

	log.Printf("Found %d pending campaigns to process", len(campaigns))

	for _, campaign := range campaigns {
		// Get user emails if not in test mode
		var userEmails []string
		if !campaign.IsTestMode() {
			userEmails, err = s.getAllUserEmails(ctx)
			if err != nil {
				log.Printf("Error getting user emails for campaign %s: %v", campaign.ID.Hex(), err)
				continue
			}
		}

		// Process campaign
		if err := s.campaignBiz.ProcessCampaign(ctx, campaign, userEmails); err != nil {
			log.Printf("Error processing campaign %s: %v", campaign.ID.Hex(), err)
		}
	}
}

// getAllUserEmails gets all user emails from database
func (s *CampaignScheduler) getAllUserEmails(ctx context.Context) ([]string, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var emails []string
	for cursor.Next(ctx) {
		var result struct {
			Email string `bson:"email"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		if result.Email != "" {
			emails = append(emails, result.Email)
		}
	}

	return emails, nil
}
