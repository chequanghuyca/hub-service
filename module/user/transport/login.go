package transport

import (
	"hub-service/component/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary User login
// @Description Login with email and password to get access token
// @Tags users
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginAPIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/users/login [post]
func Login(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginReq model.LoginRequest

		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		biz := biz.NewUserBiz(appCtx)
		loginResp, err := biz.Login(c.Request.Context(), loginReq.Email, loginReq.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.LoginAPIResponse{Status: "success", Data: *loginResp})
	}
}
