package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/challenge/biz"
	_ "hub-service/module/challenge/model"
	"hub-service/module/challenge/storage"
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
// @Success 200 {object} common.Response{data=model.Challenge} "Success"
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

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
