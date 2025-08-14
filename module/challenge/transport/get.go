package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/challenge/biz"
	"hub-service/module/challenge/model"
	"hub-service/module/challenge/storage"
	scoremodel "hub-service/module/score/model"
	scorestorage "hub-service/module/score/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetChallenge godoc
// @Summary Get a challenge by ID
// @Description Retrieve the details of a specific translation challenge by its unique ID. All authenticated users can access this endpoint.
// @Tags challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Challenge ID (MongoDB ObjectID)"
// @Success 200 {object} common.Response{data=model.ChallengeDetail} "Success"
// @Failure 400 {object} common.AppError "Invalid ID format"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Challenge not found"
// @Router /api/challenges/{id} [get]
func GetChallenge(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewGetChallengeBiz(store)

		data, err := business.GetChallenge(c.Request.Context(), id)
		if err != nil {
			panic(err)
		}

		// Enrich with user's best score and last answer if available
		var userID primitive.ObjectID
		if v, exists := c.Get("user_id"); exists {
			if uid, ok := v.(primitive.ObjectID); ok {
				userID = uid
			}
		}

		// Default: no user-specific fields
		detail := model.ChallengeDetail{Challenge: *data}

		if userID != primitive.NilObjectID {
			sStore := scorestorage.NewStorage(appCtx.GetDatabase())
			score, err := sStore.GetScoreByUserAndChallenge(c.Request.Context(), userID, data.ID)
			if err != nil {
				panic(err)
			}
			if score != nil {
				// Build ChallengeScore-like response
				cs := scoremodel.ChallengeScore{
					ChallengeID:     data.ID,
					ChallengeTitle:  data.Title,
					BestScore:       score.BestScore,
					AttemptCount:    score.AttemptCount,
					LastAttemptAt:   score.UpdatedAt,
					UserTranslation: score.UserTranslation,
					Feedback:        score.Feedback,
					Errors:          score.Errors,
					Suggestions:     score.Suggestions,
					OriginalContent: score.OriginalContent,
				}
				// Optional fields as requested
				if score.BestScore > 0 {
					b := score.BestScore
					detail.UserBestScore = &b
				}
				detail.UserScore = &cs
			}
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(detail))
	}
}
