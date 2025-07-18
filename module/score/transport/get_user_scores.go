package transport

import (
	"errors"
	"hub-service/common"
	"hub-service/core/appctx"
	challengestorage "hub-service/module/challenge/storage"
	scorebiz "hub-service/module/score/biz"
	scoremodel "hub-service/module/score/model"
	"hub-service/module/score/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserScores godoc
// @Summary Get user's scores for all challenges
// @Description Get detailed scores and summary for a specific user
// @Tags scores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} scoremodel.GetUserScoresAPIResponse
// @Failure 400 {object} common.AppError
// @Failure 401 {object} common.AppError
// @Failure 500 {object} common.AppError
// @Router /api/scores/user/{user_id} [get]
func GetUserScores(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scoremodel.GetUserScoresRequest
		if err := c.ShouldBindUri(&req); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		userID, err := primitive.ObjectIDFromHex(req.UserID)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		challengeStore := challengestorage.NewStorage(appCtx.GetDatabase())
		geminiAPIKey := appCtx.GetEnv("GEMINI_API_KEY")
		if geminiAPIKey == "" {
			panic(common.NewErrorResponse(err, "Gemini API key not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}
		geminiBaseURL := appCtx.GetEnv("GEMINI_BASE_URL")
		if geminiBaseURL == "" {
			panic(common.NewErrorResponse(errors.New("gemini base url not configured"), "Gemini base url not configured", "GEMINI_CONFIG_ERROR", "CONFIG_ERROR"))
		}
		geminiBiz := scorebiz.NewGeminiBiz(geminiAPIKey, geminiBaseURL)
		business := scorebiz.NewScoreBiz(store, challengeStore, geminiBiz)

		result, err := business.GetUserScores(c.Request.Context(), userID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
