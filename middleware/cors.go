package middleware

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CorsConnect() gin.HandlerFunc {
	godotenv.Load()

	corsConfig := cors.Config{
		AllowOrigins:     []string{os.Getenv("BASE_URL_PORTFOLIO")},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}, // Các phương thức được phép
		AllowHeaders:     []string{"Content-Type", "Authorization"},             // Các header được phép
		ExposeHeaders:    []string{"Content-Length"},                            // Header trả về mà FE có thể đọc
		AllowCredentials: true,                                                  // Cho phép gửi cookie nếu cần
		MaxAge:           12 * 60 * 60,
	}

	return cors.New(corsConfig)
}
