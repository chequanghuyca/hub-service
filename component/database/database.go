package database

import (
	"hub-service/component/mongodb"
	"log"
	"os"
)

type Database struct {
	MongoDB *mongodb.MongoDB
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

	log.Println("All database connections established successfully!")

	return &Database{
		MongoDB: mongoClient,
	}, nil
}

// Close closes all database connections
func (db *Database) Close() error {
	if db.MongoDB != nil {
		return db.MongoDB.Close()
	}
	return nil
}

// HealthCheck checks health of all database connections
func (db *Database) HealthCheck() error {
	if db.MongoDB != nil {
		return db.MongoDB.HealthCheck()
	}
	return nil
}
