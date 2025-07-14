package storage

import (
	"context"
	"hub-service/module/challenge/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) Get(ctx context.Context, id primitive.ObjectID) (*model.Challenge, error) {
	var data model.Challenge
	collection := s.db.MongoDB.GetCollection(model.CollectionName)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Storage) GetChallenge(ctx context.Context, id primitive.ObjectID) (*model.Challenge, error) {
	var data model.Challenge

	collection := s.db.MongoDB.GetCollection(model.CollectionName)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
