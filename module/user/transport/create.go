package transport

import (
	"hub-service/component/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with email, password and name
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserCreate true "User information"
// @Success 200 {object} model.CreateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/users [post]
func CreateUser(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCreate model.UserCreate

		if err := c.ShouldBindJSON(&userCreate); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		biz := biz.NewUserBiz(appCtx)
		user, err := biz.CreateUser(c.Request.Context(), &userCreate)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.CreateUserResponse{Status: "success", Data: *user})
	}
}
