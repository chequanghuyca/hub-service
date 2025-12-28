package main

import (
	"hub-service/core/appctx"
	"hub-service/docs"
	"hub-service/infrastructure/database/database"
	"hub-service/middleware"
	emailConsumer "hub-service/module/email/consumer"
	emailRepository "hub-service/module/email/repository"
	"hub-service/module/email/scheduler"
	emailSender "hub-service/module/email/sender"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Hub Service API
// @version 1.0
// @description This is a sample hub service API. You can enter your access token directly without 'Bearer ' prefix.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your access token directly (without 'Bearer ' prefix). The system will automatically handle both formats.
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

	// Start email consumer if Kafka is configured
	if appContext.GetKafka() != nil {
		emailRepo := emailRepository.NewEmailRepository(db.MongoDB.Database)
		sender := emailSender.NewSMTPSender()
		consumer := emailConsumer.NewEmailConsumer(
			appContext.GetKafka(),
			appContext.GetRedis(),
			emailRepo,
			sender,
		)
		consumer.Start()
		defer consumer.Stop()

		// Start campaign scheduler
		campaignScheduler := scheduler.NewCampaignScheduler(appContext)
		campaignScheduler.Start()
		defer campaignScheduler.Stop()
	}

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
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down gracefully...")
		os.Exit(0)
	}()

	r.Run()
}

