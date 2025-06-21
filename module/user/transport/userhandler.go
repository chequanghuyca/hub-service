package transport

import (
	"hub-service/component/appctx"
	"hub-service/module/user/biz"
	"hub-service/module/user/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userHandler struct {
	appCtx appctx.AppContext
	biz    *biz.UserBiz
}

func NewUserHandler(appCtx appctx.AppContext) *userHandler {
	return &userHandler{
		appCtx: appCtx,
		biz:    biz.NewUserBiz(appCtx),
	}
}

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
func (h *userHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userCreate model.UserCreate

		if err := c.ShouldBindJSON(&userCreate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := h.biz.CreateUser(c.Request.Context(), &userCreate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}

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
func (h *userHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginReq model.LoginRequest

		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		loginResp, err := h.biz.Login(c.Request.Context(), loginReq.Email, loginReq.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": loginResp,
		})
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.GetUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/users/{id} [get]
func (h *userHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		user, err := h.biz.GetUserByID(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body model.UserUpdate true "User update information"
// @Security ApiKeyAuth
// @Success 200 {object} model.UpdateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/users/{id} [put]
func (h *userHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		var userUpdate model.UserUpdate

		if err := c.ShouldBindJSON(&userUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := h.biz.UpdateUser(c.Request.Context(), objectID, &userUpdate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.DeleteUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /api/users/{id} [delete]
func (h *userHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		err = h.biz.DeleteUser(c.Request.Context(), objectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted successfully",
		})
	}
}

// ListUsers godoc
// @Summary List users
// @Description Get list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of users" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Security ApiKeyAuth
// @Success 200 {object} model.ListUsersResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /api/users [get]
func (h *userHandler) ListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "10")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid limit parameter",
			})
			return
		}

		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid offset parameter",
			})
			return
		}

		users, err := h.biz.ListUsers(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": users,
		})
	}
}
