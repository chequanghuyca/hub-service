package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	translations := g.Group("/translations")
	{
		// Read operations - accessible by all authenticated users
		protected := translations.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/:id", GetTranslation(appCtx))
			protected.GET("/:id/progress", GetTranslationWithProgress(appCtx))
			protected.GET("/user/:user_id/scores", GetUserTranslationScores(appCtx))
		}

		// Create and update operations - accessible by admin and super_admin
		adminProtected := translations.Group("/")
		adminProtected.Use(auth.AuthMiddleware(appCtx))
		adminProtected.Use(auth.RequireRoles(common.RoleAdmin, common.RoleSuperAdmin))
		{
			adminProtected.POST("/create", CreateTranslation(appCtx))
		}

		// User operations - accessible by all authenticated users
		userProtected := translations.Group("/")
		userProtected.Use(auth.AuthMiddleware(appCtx))
		{
			userProtected.POST("/:id/sentences/:sentence_index/translate", SubmitSentenceTranslation(appCtx))
		}
	}
}
