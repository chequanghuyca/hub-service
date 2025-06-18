package main

import (
	"hub-service/component/appctx"
	"hub-service/middleware"
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

	appContext := appctx.NewAppContext(os.Getenv("SYSTEM_SECRET_KEY"))

	r := gin.Default()

	r.Use(middleware.Recover(appContext))
	r.Use(middleware.CorsConnect())

	r.Static("/static", "./static")

	middleware.ApiServices(appContext, r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}
