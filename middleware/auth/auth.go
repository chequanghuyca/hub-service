package auth

import (
	"errors"
	"hub-service/core/appctx"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and extracts user information
func AuthMiddleware(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		tokenProvider := appCtx.GetTokenProvider()
		payload, err := tokenProvider.Validate(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", payload.UserID)
		c.Set("user_role", payload.Role)

		c.Next()
	}
}

func extractTokenFromHeader(s string) (string, error) {
	if s == "" {
		return "", errors.New("authorization header is required")
	}

	if strings.HasPrefix(s, "Bearer ") {
		token := strings.TrimPrefix(s, "Bearer ")
		if token == "" {
			return "", errors.New("token is empty after removing 'Bearer ' prefix")
		}
		return token, nil
	} else {
		return s, nil
	}
}

// RequireRoles middleware checks if user has one of the required roles
func RequireRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found in token"})
			c.Abort()
			return
		}

		hasPermission := false
		for _, role := range requiredRoles {
			if userRole == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
