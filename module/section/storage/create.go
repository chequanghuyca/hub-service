package storage

import (
	"context"
	"hub-service/module/section/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) Create(ctx context.Context, data *model.SectionCreate) error {
	return s.CreateSection(ctx, data)
}

func (s *Storage) CreateSection(ctx context.Context, data *model.SectionCreate) error {
	now := time.Now()
	data.ID = primitive.NewObjectID()
	data.CreatedAt = &now
	data.UpdatedAt = &now

	collection := s.db.MongoDB.GetCollection(model.SectionName)
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
