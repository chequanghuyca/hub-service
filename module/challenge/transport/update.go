package transport

import (
	"hub-service/common"
	"hub-service/component/appctx"
	"hub-service/module/challenge/biz"
	"hub-service/module/challenge/model"
	"hub-service/module/challenge/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateChallenge godoc
// @Summary Update a challenge
// @Description Update the details of an existing translation challenge by its ID.
// @Tags challenges
// @Accept json
// @Produce json
// @Param id path string true "Challenge ID (MongoDB ObjectID)"
// @Param challenge body model.ChallengeUpdate true "Challenge data to update"
// @Success 200 {object} common.Response{data=boolean} "Success"
// @Failure 400 {object} common.AppError "Bad request or invalid ID format"
// @Failure 404 {object} common.AppError "Challenge not found"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/challenges/{id} [patch]
func UpdateChallenge(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data model.ChallengeUpdate
		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewUpdateChallengeBiz(store)

		if err := business.UpdateChallenge(c.Request.Context(), id, &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
