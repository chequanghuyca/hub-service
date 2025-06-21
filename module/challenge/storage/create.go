package storage

import (
	"context"
	"hub-service/module/challenge/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) Create(ctx context.Context, data *model.ChallengeCreate) error {
	now := time.Now()
	data.ID = primitive.NewObjectID()
	data.CreatedAt = &now
	data.UpdatedAt = &now

	collection := s.db.MongoDB.GetCollection(model.CollectionName)
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
