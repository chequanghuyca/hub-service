package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	challengestorage "hub-service/module/challenge/storage"
	scorebiz "hub-service/module/score/biz"
	scoremodel "hub-service/module/score/model"
	"hub-service/module/score/storage"
	translatebiz "hub-service/module/translate/biz"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTotalScore godoc
// @Summary Get user's total score
// @Description Get total score, average score, and best score for a user
// @Tags scores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} scoremodel.GetTotalScoreAPIResponse
// @Failure 400 {object} common.AppError
// @Failure 401 {object} common.AppError
// @Failure 500 {object} common.AppError
// @Router /api/scores/total/{user_id} [get]
func GetTotalScore(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req scoremodel.GetTotalScoreRequest
		if err := c.ShouldBindUri(&req); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		userID, err := primitive.ObjectIDFromHex(req.UserID)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		challengeStore := challengestorage.NewStorage(appCtx.GetDatabase())
		translateBiz := translatebiz.NewTranslateBiz(appCtx, challengeStore)
		business := scorebiz.NewScoreBiz(store, challengeStore, translateBiz)

		result, err := business.GetTotalScore(c.Request.Context(), userID)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
