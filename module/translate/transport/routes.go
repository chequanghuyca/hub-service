package transport

import (
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	translate := g.Group("/translate")
	{
		// The /score endpoint now requires authentication to save the user's score.
		translate.POST("/score", auth.AuthMiddleware(appCtx), ScoreTranslationHandler(appCtx))
	}
}
