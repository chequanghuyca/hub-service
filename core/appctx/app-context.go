package appctx

import (
	"hub-service/core/auth/tokenprovider"
	"hub-service/core/auth/tokenprovider/jwt"
	"hub-service/infrastructure/database/database"
	"hub-service/infrastructure/database/redis"
	"hub-service/infrastructure/external/deepl"
	"hub-service/infrastructure/messaging/kafka"
	"log"
	"os"

	deeplgo "github.com/bounoable/deepl"
)

type AppContext interface {
	GetSecretKey() string
	GetDatabase() *database.Database
	GetTokenProvider() tokenprovider.Provider
	GetDeeplClient() *deeplgo.Client
	GetKafka() *kafka.KafkaClient
	GetRedis() *redis.RedisClient
	GetEnv(key string) string
}

type appContext struct {
	secretKey     string
	db            *database.Database
	tokenProvider tokenprovider.Provider
	deeplClient   *deeplgo.Client
	kafkaClient   *kafka.KafkaClient
}

func NewAppContext(secretKey string, db *database.Database) *appContext {
	tokenProvider := jwt.NewProvider(secretKey)
	deeplClient := deepl.NewClient()

	// Initialize Kafka (optional)
	var kafkaClient *kafka.KafkaClient
	if os.Getenv("KAFKA_BROKERS") != "" {
		var err error
		kafkaClient, err = kafka.NewKafkaClient()
		if err != nil {
			log.Printf("Warning: Failed to connect to Kafka: %v", err)
		}
	}

	return &appContext{
		secretKey:     secretKey,
		db:            db,
		tokenProvider: tokenProvider,
		deeplClient:   deeplClient,
		kafkaClient:   kafkaClient,
	}
}

func (ctx *appContext) GetSecretKey() string {
	return ctx.secretKey
}

func (ctx *appContext) GetDatabase() *database.Database {
	return ctx.db
}

func (ctx *appContext) GetTokenProvider() tokenprovider.Provider {
	return ctx.tokenProvider
}

func (ctx *appContext) GetDeeplClient() *deeplgo.Client {
	return ctx.deeplClient
}

func (ctx *appContext) GetKafka() *kafka.KafkaClient {
	return ctx.kafkaClient
}

func (ctx *appContext) GetRedis() *redis.RedisClient {
	if ctx.db != nil {
		return ctx.db.Redis
	}
	return nil
}

func (appCtx *appContext) GetEnv(key string) string {
	return os.Getenv(key)
}

