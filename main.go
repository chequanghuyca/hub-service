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
func main() {
	godotenv.Load()

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	appContext := appctx.NewAppContext(os.Getenv("SYSTEM_SECRET_KEY"), db)

	r := gin.Default()

	r.Use(middleware.Recover(appContext))
	r.Use(middleware.CorsConnect())

	r.Static("/static", "./static")

	r.GET("/health", middleware.HealthCheck(appContext))

	middleware.ApiServices(appContext, r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
