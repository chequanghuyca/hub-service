package transport

import (
	"hub-service/core/appctx"
	"hub-service/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.RouterGroup, appCtx appctx.AppContext) {
	upload := g.Group("/upload")
	upload.Use(auth.AuthMiddleware(appCtx))
	{
		upload.POST("/r2-image", UploadR2Image(appCtx))
	}
}
