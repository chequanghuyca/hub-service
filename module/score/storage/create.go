package storage

import (
	"context"
	"hub-service/module/score/model"
)

func (s *Storage) CreateScore(ctx context.Context, data *model.ScoreCreate) error {
	collection := s.db.MongoDB.GetCollection(model.CollectionName)

	_, err := collection.InsertOne(ctx, data)
	return err
}
