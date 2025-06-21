package middleware

import (
	"hub-service/component/appctx"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check requests
func HealthCheck(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		mongoErr := appCtx.GetDatabase().HealthCheck()

		if mongoErr != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"message": "Database connection failed",
				"error":   mongoErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "All services are running",
			"mongodb": "connected",
		})
	}
}
