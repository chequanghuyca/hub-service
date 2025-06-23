package storage

import (
	"context"
	"hub-service/core/appctx"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStorage struct {
	appCtx appctx.AppContext
}

func NewUserStorage(appCtx appctx.AppContext) *UserStorage {
	return &UserStorage{appCtx: appCtx}
}

func (s *UserStorage) UpdateRole(ctx context.Context, id primitive.ObjectID, role string) error {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")
	update := bson.M{"$set": bson.M{"role": role}}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
