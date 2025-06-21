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
		AllowOrigins:     []string{os.Getenv("BASE_URL_PORTFOLIO"), os.Getenv("BASE_URL_LOCAL")},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}

	return cors.New(corsConfig)
}
