package transport

import (
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	translate := g.Group("/translate")
	{
		// Protected routes - require authentication
		protected := translate.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			// The endpoint to get a challenge has been moved to the challenge module
			// challenges.GET("/:id", GetChallenge(appCtx))
			protected.POST("/score", ScoreTranslationHandler(appCtx))
		}
	}
}
