package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/challenge/biz"
	"hub-service/module/challenge/storage"
	"hub-service/module/upload/service"
	"hub-service/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteChallenge godoc
// @Summary Delete a challenge
// @Description Delete a translation challenge by its unique ID. Only admin and super_admin can access this endpoint.
// @Tags challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Challenge ID (MongoDB ObjectID)"
// @Success 200 {object} common.Response{data=boolean} "Success"
// @Failure 400 {object} common.AppError "Invalid ID format"
// @Failure 401 {object} common.AppError "Unauthorized"
// @Failure 403 {object} common.AppError "Forbidden - Only admin and super_admin can access"
// @Failure 404 {object} common.AppError "Challenge not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/challenges/{id} [delete]
func DeleteChallenge(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())

		challenge, err := store.GetChallenge(c.Request.Context(), id)
		if err != nil {
			panic(err)
		}

		if challenge.Image != "" {
			oldFileName := helper.ExtractFileNameFromURL(challenge.Image)
			if oldFileName != "" {
				r2Service, err := service.NewR2Service()
				if err == nil {
					_ = r2Service.DeleteFile(oldFileName)
				}
			}
		}

		business := biz.NewDeleteChallengeBiz(store)

		if err := business.DeleteChallenge(c.Request.Context(), id); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
