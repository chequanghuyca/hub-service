package middleware

import (
	"hub-service/component/appctx"
	ginemail "hub-service/module/email/transport"

	"github.com/gin-gonic/gin"
)

func ApiServices(appCtx appctx.AppContext, r *gin.Engine) {
	v1 := r.Group("/api")

	v1.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	ginemail.RegisterRoutes(appCtx, v1)

}
