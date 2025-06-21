package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *UserStorage) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
