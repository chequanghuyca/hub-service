package transport

import (
	"hub-service/core/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body model.UserUpdate true "User update information"
// @Security BearerAuth
// @Success 200 {object} model.UpdateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/users/{id} [put]
func UpdateUser(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		var userUpdate model.UserUpdate

		if err := c.ShouldBindJSON(&userUpdate); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		biz := biz.NewUserBiz(appCtx)
		user, err := biz.UpdateUser(c.Request.Context(), objectID, &userUpdate)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.UpdateUserResponse{Status: "success", Data: *user})
	}
}
