package transport

import (
	"errors"
	"hub-service/common"
	"hub-service/core/appctx"
	challengestorage "hub-service/module/challenge/storage"
	scorebiz "hub-service/module/score/biz"
	scoremodel "hub-service/module/score/model"
	scorestorage "hub-service/module/score/storage"
	translatebiz "hub-service/module/translate/biz"
	translatemodel "hub-service/module/translate/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ScoreTranslationHandler godoc
// @Summary Score and save user translation for a challenge
// @Description Submits a user's translation, gets a score from DeepL, and saves the result.
// @Tags translate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body translatemodel.ScoreRequest true "Translation scoring request"
// @Success 200 {object} common.Response{data=scoremodel.SubmitScoreResponse} "Success"
// @Failure 400 {object} common.AppError "Bad request - invalid input"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Challenge not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/translate/score [post]
func ScoreTranslationHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req translatemodel.ScoreRequest
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
		scoreStore := scorestorage.NewStorage(appCtx.GetDatabase())
		challengeStore := challengestorage.NewStorage(appCtx.GetDatabase())
		translateBusiness := translatebiz.NewTranslateBiz(appCtx, challengeStore)
		scoreBusiness := scorebiz.NewScoreBiz(scoreStore, challengeStore, translateBusiness)

		// Prepare request for score submission
		submitReq := &scoremodel.SubmitScoreRequest{
			ChallengeID:     req.ChallengeID,
			UserTranslation: req.UserTranslation,
		}

		// Submit and save the score
		result, err := scoreBusiness.SubmitScore(c.Request.Context(), userID, submitReq)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
