package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	challenges := g.Group("/challenges")
	{
		// Read operations - accessible by all authenticated users (admin, super_admin, client)
		protected := challenges.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/:id", GetChallenge(appCtx))
			protected.GET("", ListChallenge(appCtx))
		}

		// Write operations - only for admin and super_admin
		adminProtected := challenges.Group("/")
		adminProtected.Use(auth.AuthMiddleware(appCtx))
		adminProtected.Use(auth.RequireRoles(common.RoleAdmin, common.RoleSuperAdmin))
		{
			adminProtected.POST("", CreateChallenge(appCtx))
			adminProtected.PATCH("/:id", UpdateChallenge(appCtx))
			adminProtected.DELETE("/:id", DeleteChallenge(appCtx))
		}
	}
}
