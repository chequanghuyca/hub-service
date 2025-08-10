package transport

import (
	"errors"
	"net/http"

	"hub-service/common"
	"hub-service/core/appctx"
	scorebiz "hub-service/module/score/biz"

	"github.com/gin-gonic/gin"
)

// AIDemoResponse represents the AI analysis payload for the demo endpoint
type AIDemoResponse struct {
	OriginalText    string   `json:"original_text"`
	UserTranslation string   `json:"user_translation"`
	TargetLanguage  string   `json:"target_language"`
	Score           float64  `json:"score"`
	Feedback        string   `json:"feedback"`
	Errors          []Error  `json:"errors"`
	Suggestions     []string `json:"suggestions"`
}

// (Removed AIDemoHandler GET as per request)

// DemoGeminiScoreRequest represents the POST body for demo scoring
type DemoGeminiScoreRequest struct {
	UserTranslation string `json:"user_translation" binding:"required"`
	TargetLanguage  string `json:"target_language" binding:"required"`
}

// AIDemoScoreHandler godoc
// @Summary AI demo scoring (no auth, no persistence)
// @Description Performs Gemini AI analysis on a fixed Vietnamese sentence using the provided user translation and target language. No auth. Does not save data.
// @Tags scores
// @Accept json
// @Produce json
// @Param request body DemoGeminiScoreRequest true "Demo scoring request"
// @Success 200 {object} common.Response{data=AIDemoResponse}
// @Failure 400 {object} common.AppError
// @Failure 500 {object} common.AppError
// @Router /api/scores/ai-demo [post]
func AIDemoScoreHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DemoGeminiScoreRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		original := "Việc học tiếng Anh là rất quan trọng. Tiếng Anh là ngôn ngữ chính của quốc tế."

		geminiAPIKey := appCtx.GetEnv("GEMINI_API_KEY")
		if geminiAPIKey == "" {
			panic(common.NewErrorResponse(errors.New("gemini api key not configured"), "Gemini API key not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}
		geminiBaseURL := appCtx.GetEnv("GEMINI_BASE_URL")
		if geminiBaseURL == "" {
			panic(common.NewErrorResponse(errors.New("gemini base url not configured"), "Gemini base url not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}

		gemini := scorebiz.NewGeminiBiz(geminiAPIKey, geminiBaseURL)
		analysis, err := gemini.AnalyzeGrammar(c.Request.Context(), original, req.UserTranslation, req.TargetLanguage)
		if err != nil {
			panic(err)
		}

		var transportErrors []Error
		if len(analysis.Errors) > 0 {
			transportErrors = make([]Error, len(analysis.Errors))
			for i, e := range analysis.Errors {
				transportErrors[i] = Error{
					Type:        e.Type,
					Description: e.Description,
					Position:    e.Position,
					Correction:  e.Correction,
				}
			}
		}

		resp := &AIDemoResponse{
			OriginalText:    original,
			UserTranslation: req.UserTranslation,
			TargetLanguage:  req.TargetLanguage,
			Score:           analysis.Score,
			Feedback:        analysis.Feedback,
			Errors:          transportErrors,
			Suggestions:     analysis.Suggestions,
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(resp))
	}
}
