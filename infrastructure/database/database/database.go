package database

import (
	"hub-service/infrastructure/database/mongodb"
	"hub-service/infrastructure/database/redis"
	"log"
	"os"
)

type Database struct {
	MongoDB *mongodb.MongoDB
	Redis   *redis.RedisClient
}

func NewDatabase() (*Database, error) {
	mongoURI := os.Getenv("SYSTEM_MONGODB_URI")
	mongoDatabase := os.Getenv("SYSTEM_APP_NAME")

	if mongoURI == "" {
		log.Fatal("SYSTEM_MONGODB_URI environment variable is required")
	}
	if mongoDatabase == "" {
		log.Fatal("SYSTEM_APP_NAME environment variable is required")
	}

	mongoClient, err := mongodb.NewMongoDB(mongoURI, mongoDatabase)
	if err != nil {
		return nil, err
	}

	// Initialize Redis (optional - won't fail if not configured)
	var redisClient *redis.RedisClient
	if os.Getenv("REDIS_HOST") != "" {
		redisClient, err = redis.NewRedisClient()
		if err != nil {
			log.Printf("Warning: Failed to connect to Redis: %v", err)
			// Don't fail - Redis is optional
		}
	}

	log.Println("All database connections established successfully!")

	return &Database{
		MongoDB: mongoClient,
		Redis:   redisClient,
	}, nil
}

// Close closes all database connections
func (db *Database) Close() error {
	if db.MongoDB != nil {
		if err := db.MongoDB.Close(); err != nil {
			return err
		}
	}
	if db.Redis != nil {
		if err := db.Redis.Close(); err != nil {
			return err
		}
	}
	return nil
}

// HealthCheck checks health of all database connections
func (db *Database) HealthCheck() error {
	if db.MongoDB != nil {
		if err := db.MongoDB.HealthCheck(); err != nil {
			return err
		}
	}
	if db.Redis != nil {
		if err := db.Redis.HealthCheck(); err != nil {
			return err
		}
	}
	return nil
}
