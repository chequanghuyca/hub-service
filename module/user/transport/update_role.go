package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update role of a user (super_admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param body body model.UpdateRoleRequest true "User email and new role"
// @Security BearerAuth
// @Success 200 {object} model.UpdateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Router /api/users/set-role [patch]
func UpdateUserRole(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
			Role  string `json:"role" binding:"required,oneof=admin client super_admin"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		userRole, _ := c.Get("user_role")
		if userRole != common.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "Insufficient permissions"})
			return
		}

		biz := biz.NewUserBiz(appCtx)
		user, err := biz.GetUserByEmail(c.Request.Context(), req.Email)
		if err != nil || user == nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "User not found"})
			return
		}

		if err := biz.StoreUpdateRole(c.Request.Context(), user.ID, req.Role); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		user.Role = req.Role
		c.JSON(http.StatusOK, model.UpdateUserResponse{Status: "success", Data: model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Avatar:    user.Avatar,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}})
	}
}
