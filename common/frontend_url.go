package common

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// IsOriginAllowed checks if the origin is in the list of allowed origins
func IsOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

// GetFrontendURL determines the frontend URL from request headers or environment
func GetFrontendURL(c *gin.Context) string {
	// Get allowed origins from env (comma-separated)
	allowedOriginsEnv := os.Getenv("CORS_ALLOW_ORIGINS")

	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		// Trim spaces from allowed origins
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	// Try to get origin from Referer header first
	referer := c.GetHeader("Referer")
	if referer != "" {
		if parsedURL, err := url.Parse(referer); err == nil {
			origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
			// If CORS_ALLOW_ORIGINS is set, validate. Otherwise, trust referer
			if len(allowedOrigins) == 0 || IsOriginAllowed(origin, allowedOrigins) {
				return origin
			}
		}
	}

	// Try Origin header
	origin := c.GetHeader("Origin")
	if origin != "" {
		// If CORS_ALLOW_ORIGINS is set, validate. Otherwise, trust origin
		if len(allowedOrigins) == 0 || IsOriginAllowed(origin, allowedOrigins) {
			return origin
		}
	}

	// Fallback to first allowed origin if available
	if len(allowedOrigins) > 0 {
		return allowedOrigins[0]
	}

	// Fallback to BASE_URL_TRANSMASTER_PROD
	if baseURL := os.Getenv("BASE_URL_TRANSMASTER_PROD"); baseURL != "" {
		return baseURL
	}

	// Fallback to default localhost for development
	return "http://localhost:3000"
}
