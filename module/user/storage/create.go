package storage

import (
	"context"
	"hub-service/module/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *UserStorage) Create(ctx context.Context, userCreate *model.UserCreate) (*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	now := time.Now()
	user := &model.User{
		Email:        userCreate.Email,
		Name:         userCreate.Name,
		Avatar:       userCreate.Avatar,
		Role:         userCreate.Role,
		Provider:     userCreate.Provider,
		ProviderID:   userCreate.ProviderID,
		IsFirstLogin: true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}
