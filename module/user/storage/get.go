package storage

import (
	"context"
	"hub-service/module/user/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *UserStorage) GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	var user model.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
