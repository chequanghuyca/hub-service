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
		users.POST("/", CreateUser(appCtx))
		users.POST("/login", Login(appCtx))
		users.POST("/refresh", RefreshToken(appCtx))

		protected := users.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/:id", GetUserByID(appCtx))
			protected.PUT("/:id", UpdateUser(appCtx))

			// List users: only for super admin and admin
			protected.GET("/", auth.RequireRoles(common.RoleSuperAdmin, common.RoleAdmin), ListUsers(appCtx))

			// Delete user: only for super admin
			protected.DELETE("/:id", auth.RequireRoles(common.RoleSuperAdmin), DeleteUser(appCtx))
		}
	}
}
