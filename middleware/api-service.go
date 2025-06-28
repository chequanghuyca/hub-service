package middleware

import (
	"hub-service/core/appctx"
	challengeTransport "hub-service/module/challenge/transport"
	ginemail "hub-service/module/email/transport"
	scoreTransport "hub-service/module/score/transport"
	ginuser "hub-service/module/user/transport"

	"github.com/gin-gonic/gin"
)

func ApiServices(appCtx appctx.AppContext, r *gin.Engine) {
	v1 := r.Group("/api")

	v1.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello!",
		})
	})

	ginuser.RegisterRoutes(appCtx, v1)
	ginemail.RegisterRoutes(appCtx, v1)
	challengeTransport.RegisterRoutes(v1, appCtx)
	scoreTransport.RegisterRoutes(v1, appCtx)
}
