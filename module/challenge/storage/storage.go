package storage

import "hub-service/infrastructure/database/database"

type Storage struct {
	db *database.Database
}

func NewStorage(db *database.Database) *Storage {
	return &Storage{db: db}
}
