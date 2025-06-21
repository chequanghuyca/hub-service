package transport

import (
	"hub-service/common"
	"hub-service/component/appctx"
	"hub-service/module/challenge/biz"
	"hub-service/module/challenge/model"
	"hub-service/module/challenge/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateChallenge godoc
// @Summary Create a new challenge
// @Description Create a new translation challenge and store it in the database.
// @Tags challenges
// @Accept json
// @Produce json
// @Param challenge body model.ChallengeCreate true "Challenge data to create"
// @Success 200 {object} common.Response{data=string} "Successfully created. Returns the ID of the new challenge."
// @Failure 400 {object} common.AppError "Bad request"
// @Failure 500 {object} common.AppError "Internal server error"
// @Router /api/challenges [post]
func CreateChallenge(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.ChallengeCreate
		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := storage.NewStorage(appCtx.GetDatabase())
		business := biz.NewCreateChallengeBiz(store)

		if err := business.CreateChallenge(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.ID))
	}
}
