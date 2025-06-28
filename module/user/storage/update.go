package storage

import (
	"context"
	"hub-service/module/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *UserStorage) Update(ctx context.Context, id primitive.ObjectID, userUpdate *model.UserUpdate) (*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	updateData := bson.M{}

	if userUpdate.Name != "" {
		updateData["name"] = userUpdate.Name
	}

	if len(updateData) == 0 {
		return s.GetByID(ctx, id)
	}

	updateData["updated_at"] = time.Now()

	update := bson.M{
		"$set": updateData,
	}

	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user model.User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLoginStatus updates user's login status and returns if this is the first login
func (s *UserStorage) UpdateLoginStatus(ctx context.Context, id primitive.ObjectID) (bool, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	now := time.Now()

	// First, get current user to check if it's first login
	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return false, err
	}

	isFirstLogin := user.IsFirstLogin

	// Update login status
	updateData := bson.M{
		"last_login_at": now,
		"updated_at":    now,
	}

	// If this is the first login, update the flag
	if isFirstLogin {
		updateData["is_first_login"] = false
	}

	update := bson.M{
		"$set": updateData,
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)

	return isFirstLogin, err
}
