package transport

import (
	"encoding/json"
	"errors"
	"hub-service/common"
	"hub-service/core/appctx"
	challengestorage "hub-service/module/challenge/storage"
	scorebiz "hub-service/module/score/biz"
	scoremodel "hub-service/module/score/model"
	"hub-service/module/score/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GeminiScoreRequest represents the request for Gemini scoring
type GeminiScoreRequest struct {
	ChallengeID     primitive.ObjectID `json:"challenge_id" binding:"required"`
	UserTranslation string             `json:"user_translation" binding:"required"`
	TargetLanguage  string             `json:"target_language" binding:"required"`
}

// GeminiScoreResponse represents the enhanced response with Gemini analysis
type GeminiScoreResponse struct {
	Score       float64            `json:"score"`
	Feedback    string             `json:"feedback"`
	Errors      []Error            `json:"errors"`
	Suggestions []string           `json:"suggestions"`
	ChallengeID primitive.ObjectID `json:"challenge_id"`
	UserID      primitive.ObjectID `json:"user_id"`
	CreatedAt   int64              `json:"created_at"`
}

type Error struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	Correction  string `json:"correction"`
}

// GeminiTranslateRequest represents the request for Gemini translation
// (not scoring, just translation)
type GeminiTranslateRequest struct {
	Text           string `json:"text" binding:"required"`
	SourceLanguage string `json:"source_language" binding:"required"`
	TargetLanguage string `json:"target_language" binding:"required"`
}

type GeminiTranslateResponse struct {
	OriginalText   string  `json:"original_text"`
	TranslatedText string  `json:"translated_text"`
	SourceLanguage string  `json:"source_language"`
	TargetLanguage string  `json:"target_language"`
	Confidence     float64 `json:"confidence"`
	Explanation    string  `json:"explanation"`
}

// GeminiScoreHandler godoc
// @Summary Score and analyze grammar using Gemini AI
// @Description Analyzes user translation using Gemini AI for grammar, syntax, and language accuracy
// @Tags scores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GeminiScoreRequest true "Gemini scoring request"
// @Success 200 {object} common.Response{data=GeminiScoreResponse} "Success"
// @Failure 400 {object} common.AppError "Bad request - invalid input"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Challenge not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/scores/ai-translate [post]
func GeminiScoreHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GeminiScoreRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		// Get user ID from context (it's set by the auth middleware)
		userIDVal, exists := c.Get("user_id")
		if !exists {
			panic(common.NewUnauthorized(errors.New("user not authenticated"), "user not authenticated", "USER_NOT_FOUND"))
		}

		userID, ok := userIDVal.(primitive.ObjectID)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("user_id in context is not of type ObjectID")))
		}

		// Create dependencies for score biz
		scoreStore := storage.NewStorage(appCtx.GetDatabase())
		challengeStore := challengestorage.NewStorage(appCtx.GetDatabase())

		geminiAPIKey := appCtx.GetEnv("GEMINI_API_KEY")
		if geminiAPIKey == "" {
			panic(common.NewErrorResponse(errors.New("gemini api key not configured"), "Gemini API key not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}

		geminiBaseURL := appCtx.GetEnv("GEMINI_BASE_URL")
		if geminiBaseURL == "" {
			panic(common.NewErrorResponse(errors.New("gemini base url not configured"), "Gemini base url not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}

		geminiBiz := scorebiz.NewGeminiBiz(geminiAPIKey, geminiBaseURL)
		business := scorebiz.NewScoreBiz(scoreStore, challengeStore, geminiBiz)

		// Convert request to SubmitScoreRequest format
		submitReq := &scoremodel.SubmitScoreRequest{
			ChallengeID:     req.ChallengeID.Hex(),
			UserTranslation: req.UserTranslation,
		}

		// Use ScoreBiz to analyze and save to database
		result, err := business.SubmitScore(c.Request.Context(), userID, submitReq, req.TargetLanguage)
		if err != nil {
			panic(err)
		}

		// Convert result back to GeminiScoreResponse format
		var modelErrors []Error
		if result.Errors != "" {
			// Parse errors JSON string back to slice
			var errors []scorebiz.Error
			if err := json.Unmarshal([]byte(result.Errors), &errors); err == nil {
				modelErrors = make([]Error, len(errors))
				for i, err := range errors {
					modelErrors[i] = Error{
						Type:        err.Type,
						Description: err.Description,
						Position:    err.Position,
						Correction:  err.Correction,
					}
				}
			}
		}

		var suggestions []string
		if result.Suggestions != "" {
			// Parse suggestions JSON string back to slice
			if err := json.Unmarshal([]byte(result.Suggestions), &suggestions); err != nil {
				suggestions = []string{}
			}
		}

		response := &GeminiScoreResponse{
			Score:       result.Score,
			Feedback:    result.Feedback,
			Errors:      modelErrors,
			Suggestions: suggestions,
			ChallengeID: req.ChallengeID,
			UserID:      userID,
			CreatedAt:   time.Now().Unix(),
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(response))
	}
}
