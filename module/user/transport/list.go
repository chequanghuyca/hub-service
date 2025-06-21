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
// @Description Get a list of users with pagination, sorting, and search
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (minimum: 1, default: 1)" minimum(1)
// @Param limit query int false "Number of items per page (minimum: 1, maximum: 100, default: 10)" minimum(1) maximum(100)
// @Param sort_by query string false "Sort by field (name, email, created_at, updated_at)" Enums(name, email, created_at, updated_at) default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" Enums(asc, desc) default(desc)
// @Param search query string false "Search by name or email (case-insensitive)"
// @Success 200 {object} model.ListUsersResponse "Returns users list with pagination metadata"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /api/users [get]
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
		if req.SortBy == "" {
			req.SortBy = "created_at"
		}
		if req.SortOrder == "" {
			req.SortOrder = "desc"
		}

		biz := biz.NewUserBiz(appCtx)
		users, metadata, err := biz.ListUsers(
			c.Request.Context(),
			req.Page,
			req.Limit,
			req.SortBy,
			req.SortOrder,
			req.Search,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.ListUsersResponse{
			Status:   "success",
			Data:     users,
			Metadata: *metadata,
		})
	}
}
