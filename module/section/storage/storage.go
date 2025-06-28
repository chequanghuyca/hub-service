package storage

import (
	"context"
	"hub-service/infrastructure/database/database"
	"hub-service/module/section/model"
)

type Storage struct {
	db *database.Database
}

func NewStorage(db *database.Database) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Create(ctx context.Context, data *model.SectionCreate) error {
	return s.CreateSection(ctx, data)
}
