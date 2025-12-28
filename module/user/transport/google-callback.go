package transport

import (
	"fmt"
	"hub-service/common"
	"hub-service/core/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// GoogleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Processes the OAuth callback from Google and creates/logs in user
// @Tags users
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State parameter for CSRF protection"
// @Success 302 {string} string "Redirect to frontend with token"
// @Failure 400 {object} model.ErrorResponse
// @Router /api/users/auth/google/callback [get]
func GoogleCallback(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get code and state from query params
		code := c.Query("code")
		state := c.Query("state")

		if code == "" || state == "" {
			redirectToFrontendWithError(c, "missing_parameters", "Missing code or state parameter")
			return
		}

		// Validate state and get frontend URL
		stateManager := appCtx.GetStateManager()
		isValid, frontendURL := stateManager.Validate(state)
		if !isValid {
			redirectToFrontendWithError(c, "invalid_state", "Invalid or expired state parameter")
			return
		}

		// Use frontend URL from state if available, otherwise fallback to getFrontendURL
		if frontendURL == "" {
			frontendURL = common.GetFrontendURL(c)
		}

		// Exchange code for token
		googleOAuth := appCtx.GetGoogleOAuth()
		token, err := googleOAuth.ExchangeCode(c.Request.Context(), code)
		if err != nil {
			redirectToFrontendWithError(c, "exchange_failed", "Failed to exchange authorization code")
			return
		}

		// Get user info from Google
		userInfo, err := googleOAuth.GetUserInfo(c.Request.Context(), token.AccessToken)
		if err != nil {
			redirectToFrontendWithError(c, "userinfo_failed", "Failed to get user information from Google")
			return
		}

		// Create or update user in database
		userBiz := biz.NewUserBiz(appCtx)
		loginReq := &model.SocialLoginRequest{
			Email:      userInfo.Email,
			Name:       userInfo.Name,
			Avatar:     userInfo.Picture,
			Provider:   "google",
			ProviderID: userInfo.ID,
			IdToken:    token.AccessToken, // We use the access token here since we already validated with Google
		}

		loginResp, err := userBiz.SocialLoginOAuth(c.Request.Context(), loginReq)
		if err != nil {
			redirectToFrontendWithError(c, "login_failed", fmt.Sprintf("Failed to create or login user: %v", err))
			return
		}

		// Redirect to frontend with access token and refresh token
		redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&refresh_token=%s", frontendURL, loginResp.AccessToken, loginResp.RefreshToken)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}

// Helper function to redirect to frontend with error
func redirectToFrontendWithError(c *gin.Context, errorCode, errorMessage string) {
	frontendURL := common.GetFrontendURL(c)
	redirectURL := fmt.Sprintf("%s/login?error=%s&message=%s", frontendURL, errorCode, url.QueryEscape(errorMessage))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
