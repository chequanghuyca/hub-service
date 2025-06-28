package transport

import (
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
		}
	}
}
