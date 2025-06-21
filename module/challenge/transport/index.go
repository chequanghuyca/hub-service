package transport

import (
	"hub-service/component/appctx"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	challenges := g.Group("/challenges")
	{
		challenges.POST("", CreateChallenge(appCtx))
		challenges.GET("/:id", GetChallenge(appCtx))
		challenges.GET("", ListChallenge(appCtx))
		challenges.PATCH("/:id", UpdateChallenge(appCtx))
		challenges.DELETE("/:id", DeleteChallenge(appCtx))
	}
}
