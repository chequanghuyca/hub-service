package middleware

import (
	"hub-service/common"
	"hub-service/core/appctx"

	"github.com/gin-gonic/gin"
)

func Recover(ac appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					return
				}

				appErr := common.ErrInvalidRequest(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				return
			}
		}()

		c.Next()
	}
}
