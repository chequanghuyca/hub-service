package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/email/biz"
	"hub-service/module/email/model"
	"hub-service/module/email/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailHandler struct {
	appCtx      appctx.AppContext
	business    biz.EmailBusiness
	campaignBiz biz.CampaignBusiness
}

func NewEmailHandler(appCtx appctx.AppContext) *EmailHandler {
	db := appCtx.GetDatabase().MongoDB.Database
	repo := repository.NewEmailRepository(db)
	campaignRepo := repository.NewCampaignRepository(db)
	business := biz.NewEmailBusiness(repo, appCtx.GetKafka(), appCtx.GetRedis())
	campaignBiz := biz.NewCampaignBusiness(campaignRepo, business)

	return &EmailHandler{
		appCtx:      appCtx,
		business:    business,
		campaignBiz: campaignBiz,
	}
}

// RegisterRoutes registers email routes
func RegisterRoutes(appCtx appctx.AppContext, router *gin.RouterGroup) {
	handler := NewEmailHandler(appCtx)

	email := router.Group("/email")
	{
		email.POST("/send", handler.SendEmail)
		email.POST("/send-bulk", handler.SendBulkEmails)
		email.GET("/logs", handler.GetEmailLogs)
		email.GET("/logs/:id", handler.GetEmailByID)

		// Campaign endpoints
		campaigns := email.Group("/campaigns")
		{
			campaigns.POST("", handler.CreateCampaign)
			campaigns.GET("", handler.ListCampaigns)
			campaigns.GET("/:id", handler.GetCampaign)
			campaigns.PUT("/:id", handler.UpdateCampaign)
			campaigns.DELETE("/:id", handler.CancelCampaign)
		}
	}
}

// SendEmail godoc
// @Summary Send a single email
// @Description Queue a single email for sending via Kafka
// @Tags Email
// @Accept json
// @Produce json
// @Param request body model.SendEmailRequest true "Email request"
// @Success 200 {object} common.SuccessResponse{data=model.EmailResponse}
// @Failure 400 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/send [post]
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req model.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid request body", "ErrInvalidRequest"))
		return
	}

	resp, err := h.business.QueueEmail(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to queue email", "ErrQueueEmail"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(resp))
}

// SendBulkEmails godoc
// @Summary Send bulk emails
// @Description Queue multiple emails for sending via Kafka
// @Tags Email
// @Accept json
// @Produce json
// @Param request body model.SendBulkEmailRequest true "Bulk email request"
// @Success 200 {object} common.SuccessResponse{data=[]model.EmailResponse}
// @Failure 400 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/send-bulk [post]
func (h *EmailHandler) SendBulkEmails(c *gin.Context) {
	var req model.SendBulkEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid request body", "ErrInvalidRequest"))
		return
	}

	responses, err := h.business.QueueBulkEmails(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to queue emails", "ErrQueueEmails"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(responses))
}

// GetEmailLogs godoc
// @Summary Get email logs
// @Description Get paginated email logs
// @Tags Email
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/logs [get]
func (h *EmailHandler) GetEmailLogs(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	emails, total, err := h.business.GetEmailLogs(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to get email logs", "ErrGetEmailLogs"))
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Data: gin.H{
			"emails": emails,
			"total":  total,
			"page":   page,
			"limit":  limit,
		},
	})
}

// GetEmailByID godoc
// @Summary Get email by ID
// @Description Get a specific email log by ID
// @Tags Email
// @Accept json
// @Produce json
// @Param id path string true "Email ID"
// @Success 200 {object} common.SuccessResponse{data=model.EmailMessage}
// @Failure 404 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/logs/{id} [get]
func (h *EmailHandler) GetEmailByID(c *gin.Context) {
	id := c.Param("id")

	email, err := h.business.GetEmailByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid email ID", "ErrInvalidEmailID"))
		return
	}

	if email == nil {
		c.JSON(http.StatusNotFound, common.NewCustomError(nil, "Email not found", "ErrEmailNotFound"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(email))
}

// CreateCampaign godoc
// @Summary Create a scheduled email campaign
// @Description Create a new email campaign to be sent at scheduled time
// @Tags Email Campaigns
// @Accept json
// @Produce json
// @Param request body model.CreateCampaignRequest true "Campaign request"
// @Success 200 {object} common.SuccessResponse{data=model.CampaignResponse}
// @Failure 400 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/campaigns [post]
func (h *EmailHandler) CreateCampaign(c *gin.Context) {
	var req model.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid request body", "ErrInvalidRequest"))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userId")
	var createdBy primitive.ObjectID
	if exists {
		if id, ok := userID.(primitive.ObjectID); ok {
			createdBy = id
		}
	}

	resp, err := h.campaignBiz.CreateCampaign(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "ErrCreateCampaign"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(resp))
}

// ListCampaigns godoc
// @Summary List email campaigns
// @Description Get paginated list of email campaigns
// @Tags Email Campaigns
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} common.SuccessResponse{data=model.CampaignListResponse}
// @Security BearerAuth
// @Router /api/email/campaigns [get]
func (h *EmailHandler) ListCampaigns(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	resp, err := h.campaignBiz.ListCampaigns(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to list campaigns", "ErrListCampaigns"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(resp))
}

// GetCampaign godoc
// @Summary Get campaign by ID
// @Description Get a specific campaign by ID
// @Tags Email Campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} common.SuccessResponse{data=model.CampaignResponse}
// @Failure 404 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/campaigns/{id} [get]
func (h *EmailHandler) GetCampaign(c *gin.Context) {
	id := c.Param("id")

	campaign, err := h.campaignBiz.GetCampaign(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid campaign ID", "ErrInvalidCampaignID"))
		return
	}

	if campaign == nil {
		c.JSON(http.StatusNotFound, common.NewCustomError(nil, "Campaign not found", "ErrCampaignNotFound"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(campaign))
}

// UpdateCampaign godoc
// @Summary Update a campaign
// @Description Update a pending campaign
// @Tags Email Campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Param request body model.UpdateCampaignRequest true "Update request"
// @Success 200 {object} common.SuccessResponse{data=model.CampaignResponse}
// @Failure 400 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/campaigns/{id} [put]
func (h *EmailHandler) UpdateCampaign(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid request body", "ErrInvalidRequest"))
		return
	}

	resp, err := h.campaignBiz.UpdateCampaign(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "ErrUpdateCampaign"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(resp))
}

// CancelCampaign godoc
// @Summary Cancel a campaign
// @Description Cancel a pending campaign
// @Tags Email Campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.SuccessResponse
// @Security BearerAuth
// @Router /api/email/campaigns/{id} [delete]
func (h *EmailHandler) CancelCampaign(c *gin.Context) {
	id := c.Param("id")

	if err := h.campaignBiz.CancelCampaign(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "ErrCancelCampaign"))
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]string{
		"message": "Campaign cancelled successfully",
	}))
}
