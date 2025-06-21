package transport

import (
	"hub-service/component/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListUsers godoc
// @Summary List users
// @Description Get a list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (minimum: 1, default: 1)" minimum(1)
// @Param limit query int false "Number of items per page (minimum: 1, maximum: 100, default: 10)" minimum(1) maximum(100)
// @Success 200 {object} model.ListUsersResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func ListUsers(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.ListUsersRequest

		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		// Set default values if not provided
		if req.Page == 0 {
			req.Page = 1
		}
		if req.Limit == 0 {
			req.Limit = 10
		}

		biz := biz.NewUserBiz(appCtx)
		users, err := biz.ListUsers(c.Request.Context(), req.Page, req.Limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.ListUsersResponse{Status: "success", Data: users})
	}
}
