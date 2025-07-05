package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	sections := g.Group("/sections")
	{
		// Read operations - accessible by all authenticated users (admin, super_admin, client)
		protected := sections.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/list", ListSection(appCtx))
			protected.GET("/:id", GetSection(appCtx))
		}

		// Update operations - accessible by admin and super_admin
		adminProtected := sections.Group("/")
		adminProtected.Use(auth.AuthMiddleware(appCtx))
		adminProtected.Use(auth.RequireRoles(common.RoleAdmin, common.RoleSuperAdmin))
		{
			adminProtected.PATCH("/:id", UpdateSection(appCtx))
			adminProtected.POST("/create", CreateSection(appCtx))
			adminProtected.DELETE("/:id", DeleteSection(appCtx))
		}
	}
}
