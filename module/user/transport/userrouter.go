package transport

import (
	"hub-service/component/appctx"
	"hub-service/component/auth"

	"github.com/gin-gonic/gin"
)

func UserRoute(appCtx appctx.AppContext, v1 *gin.RouterGroup) {
	userHandler := NewUserHandler(appCtx)

	users := v1.Group("/users")
	{
		// Public routes
		users.POST("/", userHandler.CreateUser())
		users.POST("/login", userHandler.Login())

		// Protected routes (require authentication)
		protected := users.Group("/")
		protected.Use(auth.AuthMiddleware(appCtx))
		{
			protected.GET("/", userHandler.ListUsers())
			protected.GET("/:id", userHandler.GetUserByID())
			protected.PUT("/:id", userHandler.UpdateUser())
			protected.DELETE("/:id", userHandler.DeleteUser())
		}
	}
}
