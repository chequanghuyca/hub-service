package transport

import (
	"errors"
	"hub-service/common"
	"hub-service/core/appctx"
	challengestorage "hub-service/module/challenge/storage"
	scorebiz "hub-service/module/score/biz"
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

		// Create dependencies for score biz (similar to score_translation.go)
		// Remove unused variable scoreStore
		// scoreStore := scorestorage.NewStorage(appCtx.GetDatabase())
		challengeStore := challengestorage.NewStorage(appCtx.GetDatabase())

		geminiAPIKey := appCtx.GetEnv("GEMINI_API_KEY")
		if geminiAPIKey == "" {
			panic(common.NewErrorResponse(errors.New("gemini api key not configured"), "Gemini API key not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}

		geminiBaseURL := appCtx.GetEnv("GEMINI_BASE_URL")
		if geminiBaseURL == "" {
			panic(common.NewErrorResponse(errors.New("gemini base url not configured"), "Gemini base url not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}

		geminiAnalyzer := scorebiz.NewGeminiBiz(geminiAPIKey, geminiBaseURL)

		challenge, err := challengeStore.GetChallenge(c.Request.Context(), req.ChallengeID)
		if err != nil {
			panic(common.ErrEntityNotFound("challenge", err))
		}

		analysis, err := geminiAnalyzer.AnalyzeGrammar(
			c.Request.Context(),
			challenge.Content,
			req.UserTranslation,
			req.TargetLanguage,
		)
		if err != nil {
			panic(err)
		}

		modelErrors := make([]Error, len(analysis.Errors))
		for i, err := range analysis.Errors {
			modelErrors[i] = Error{
				Type:        err.Type,
				Description: err.Description,
				Position:    err.Position,
				Correction:  err.Correction,
			}
		}

		response := &GeminiScoreResponse{
			Score:       analysis.Score,
			Feedback:    analysis.Feedback,
			Errors:      modelErrors,
			Suggestions: analysis.Suggestions,
			ChallengeID: req.ChallengeID,
			UserID:      userID,
			CreatedAt:   time.Now().Unix(),
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(response))
	}
}
