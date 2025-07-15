package transport

import (
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	scores := g.Group("/scores")
	{
		// All score operations require authentication
		protected := scores.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/user/:user_id", GetUserScores(appCtx))
			protected.POST("/ai-translate", auth.AuthMiddleware(appCtx), GeminiScoreHandler(appCtx))
		}
	}
}
