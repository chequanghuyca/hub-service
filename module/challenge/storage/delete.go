package storage

import (
	"context"
	"hub-service/module/challenge/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := s.db.MongoDB.GetCollection(model.CollectionName)
	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(ctx, filter)
	return err
}
