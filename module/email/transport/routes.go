package ginemail

import (
	"hub-service/core/appctx"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(appCtx appctx.AppContext, router *gin.RouterGroup) {
	email := router.Group("/email")

	email.POST("/single", SingleMail(appCtx))
	email.POST("/multiple", MultipleMail(appCtx))
	email.POST("/response-portfolio", ResponseEmailPortfolio(appCtx))
}
