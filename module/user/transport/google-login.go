package transport

import (
	"hub-service/common"
	"hub-service/core/appctx"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirects user to Google OAuth consent page
// @Tags users
// @Accept json
// @Produce json
// @Success 302 {string} string "Redirect to Google OAuth"
// @Router /api/users/auth/google/login [get]
func GoogleLogin(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get frontend URL from request (for redirect after OAuth)
		frontendURL := common.GetFrontendURL(c)

		// Generate state for CSRF protection with frontend URL
		stateManager := appCtx.GetStateManager()
		state, err := stateManager.Generate(frontendURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate state",
			})
			return
		}

		// Get Google OAuth provider
		googleOAuth := appCtx.GetGoogleOAuth()
		authURL := googleOAuth.GetAuthURL(state)

		// Redirect to Google OAuth
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
