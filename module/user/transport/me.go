package transport

import (
	"hub-service/core/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"hub-service/module/upload/service"
	"hub-service/utils/helper"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// GetMe godoc
// @Summary Get current user profile
// @Description Get current user's profile information using access token
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.GetUserResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/users/me [get]
func GetMe(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
			return
		}

		objectID, ok := userID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		biz := biz.NewUserBiz(appCtx)
		user, err := biz.GetUserByID(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "User not found"})
			return
		}

		c.JSON(http.StatusOK, model.GetUserResponse{Status: "success", Data: *user})
	}
}

// UpdateMe godoc
// @Summary Update current user profile
// @Description Update current user's profile information using access token
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserUpdate true "User update information"
// @Security BearerAuth
// @Success 200 {object} model.UpdateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/users/me [patch]
func UpdateMe(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
			return
		}

		objectID, ok := userID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		var userUpdate model.UserUpdate
		if err := c.ShouldBindJSON(&userUpdate); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		biz := biz.NewUserBiz(appCtx)

		// Get current user for avatar cleanup
		currentUser, err := biz.GetUserByID(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "User not found"})
			return
		}

		// Handle avatar cleanup if avatar is being updated
		if userUpdate.Avatar != "" {
			if currentUser.Avatar != "" && currentUser.Avatar != userUpdate.Avatar {
				oldFileName := helper.ExtractFileNameFromURL(currentUser.Avatar)
				if oldFileName != "" {
					r2Service, err := service.NewR2Service()
					if err == nil {
						_ = r2Service.DeleteFile(oldFileName)
					}
				}
			}
		}

		user, err := biz.UpdateUser(c.Request.Context(), objectID, &userUpdate)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.UpdateUserResponse{Status: "success", Data: *user})
	}
}
