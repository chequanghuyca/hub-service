package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	challengestorage "hub-service/module/challenge/storage"
	"hub-service/module/translate/biz"
	"hub-service/module/translate/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ScoreTranslationHandler godoc
// @Summary Score user translation
// @Description Score user's translation by comparing it with DeepL translation using Sørensen-Dice similarity
// @Tags translate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.ScoreRequest true "Translation scoring request"
// @Success 200 {object} common.Response{data=model.ScoreResponse} "Success"
// @Failure 400 {object} common.AppError "Bad request - invalid input"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 404 {object} common.AppError "Challenge not found"
// @Failure 500 {object} common.AppError "Internal server error - DeepL API error or scoring error"
// @Router /api/translate/score [post]
func ScoreTranslationHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.ScoreRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(err, "bad_request", err.Error(), "BAD_REQUEST"))
			return
		}

		// Create dependencies and inject them
		store := challengestorage.NewStorage(appCtx.GetDatabase())
		b := biz.NewTranslateBiz(appCtx, store)

		score, err := b.ScoreTranslation(c.Request.Context(), req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, common.NewErrorResponse(err, "server_error", err.Error(), "SCORE_ERROR"))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(score))
	}
}
