package transport

import (
	"hub-service/core/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token and refresh token pair using a valid refresh token.
// @Tags users
// @Accept json
// @Produce json
// @Param token body model.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} model.LoginAPIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/users/refresh [post]
func RefreshToken(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RefreshTokenRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		userBiz := biz.NewUserBiz(appCtx)
		loginResp, err := userBiz.RefreshToken(c.Request.Context(), req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.LoginAPIResponse{Status: "success", Data: *loginResp})
	}
}
