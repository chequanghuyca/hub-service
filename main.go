package main

import (
	"hub-service/component/appctx"
	"hub-service/component/database"
	"hub-service/middleware"
	"log"
	"os"

	_ "hub-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Hub Service API
// @version 1.0
// @description This is a sample hub service API.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.
func main() {
	godotenv.Load()

	// Initialize database connections
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize app context
	appContext := appctx.NewAppContext(os.Getenv("SYSTEM_SECRET_KEY"), db)

	// Setup router
	r := gin.Default()

	// Apply middlewares
	r.Use(middleware.Recover(appContext))
	r.Use(middleware.CorsConnect())

	// Static files
	r.Static("/static", "./static")

	// Health check endpoint
	r.GET("/health", middleware.HealthCheck(appContext))

	// API services
	middleware.ApiServices(appContext, r)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
