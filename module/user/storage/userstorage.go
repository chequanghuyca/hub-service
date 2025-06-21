package storage

import (
	"context"
	"hub-service/component/appctx"
	"hub-service/module/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userStorage struct {
	appCtx appctx.AppContext
}

func NewUserStorage(appCtx appctx.AppContext) *userStorage {
	return &userStorage{appCtx: appCtx}
}

func (s *userStorage) Create(ctx context.Context, userCreate *model.UserCreate) (*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	now := time.Now()
	user := &model.User{
		Email:     userCreate.Email,
		Password:  userCreate.Password, // In production, hash the password
		Name:      userCreate.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *userStorage) GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
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

func (s *userStorage) GetByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (s *userStorage) Update(ctx context.Context, id primitive.ObjectID, userUpdate *model.UserUpdate) (*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	update := bson.M{
		"$set": bson.M{
			"name":       userUpdate.Name,
			"updated_at": time.Now(),
		},
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

func (s *userStorage) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *userStorage) List(ctx context.Context, limit, offset int64) ([]*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
