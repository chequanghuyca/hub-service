package transport

import (
	"hub-service/core/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SocialLogin godoc
// @Summary Social login (Google)
// @Description Login or register user via Google OAuth
// @Tags users
// @Accept json
// @Produce json
// @Param social body model.SocialLoginRequest true "Social login info"
// @Success 200 {object} model.LoginAPIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/users/social-login [post]
func SocialLogin(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.SocialLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}
		biz := biz.NewUserBiz(appCtx)
		loginResp, err := biz.SocialLogin(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, model.LoginAPIResponse{Status: "success", Data: *loginResp})
	}
}
