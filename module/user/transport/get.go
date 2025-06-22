package transport

import (
	"hub-service/core/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} model.GetUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/users/{id} [get]
func GetUserByID(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		biz := biz.NewUserBiz(appCtx)
		user, err := biz.GetUserByID(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.GetUserResponse{Status: "success", Data: *user})
	}
}
