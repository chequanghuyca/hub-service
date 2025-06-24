package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(appCtx appctx.AppContext, router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/social-login", SocialLogin(appCtx))
		users.POST("/refresh", RefreshToken(appCtx))

		protected := users.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/:id", GetUserByID(appCtx))
			protected.GET("/", ListUsers(appCtx))
			protected.PUT("/:id", UpdateUser(appCtx))

			// Delete user: only for super admin
			protected.DELETE("/:id", auth.RequireRoles(common.RoleSuperAdmin, common.RoleAdmin), DeleteUser(appCtx))
			protected.PUT("/role", auth.RequireRoles(common.RoleSuperAdmin), UpdateUserRole(appCtx))
		}
	}
}
