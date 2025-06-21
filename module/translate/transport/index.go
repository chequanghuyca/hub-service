package transport

import (
	"hub-service/component/appctx"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	translate := g.Group("/translate")
	{
		// The endpoint to get a challenge has been moved to the challenge module
		// challenges.GET("/:id", GetChallenge(appCtx))
		translate.POST("/score", ScoreTranslationHandler(appCtx))
	}
}
