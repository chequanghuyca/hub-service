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
		users.GET("/list", ListUsers(appCtx))
		users.POST("/social-login", SocialLogin(appCtx))
		users.POST("/refresh", RefreshToken(appCtx))

		protected := users.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/me", GetMe(appCtx))
			protected.PATCH("/me", UpdateMe(appCtx))
			protected.GET("/:id", GetUserByID(appCtx))
			protected.PATCH("/:id", UpdateUser(appCtx))
			protected.DELETE("/:id", auth.RequireRoles(common.RoleSuperAdmin, common.RoleAdmin), DeleteUser(appCtx))
			protected.PATCH("/set-role", auth.RequireRoles(common.RoleSuperAdmin), UpdateUserRole(appCtx))
		}
	}
}
