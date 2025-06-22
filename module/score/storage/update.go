package storage

import (
	"context"
	"hub-service/module/score/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) UpdateScore(ctx context.Context, id primitive.ObjectID, data *model.ScoreUpdate) error {
	collection := s.db.MongoDB.GetCollection(model.CollectionName)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": data}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
