package storage

import (
	"context"
	"hub-service/module/user/model"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *UserStorage) List(ctx context.Context, limit, offset int64, sortBy, sortOrder, search string) ([]*model.User, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	filter := bson.M{}
	if search != "" {
		searchRegex := bson.M{"$regex": search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": searchRegex},
			{"email": searchRegex},
		}
	}

	// Build sort options
	sortValue := 1 // ASC by default
	if strings.ToLower(sortOrder) == "desc" {
		sortValue = -1
	}

	// Map sortBy to actual field names
	sortField := "created_at" // default
	switch strings.ToLower(sortBy) {
	case "name":
		sortField = "name"
	case "email":
		sortField = "email"
	case "created_at":
		sortField = "created_at"
	case "updated_at":
		sortField = "updated_at"
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{sortField: sortValue})

	cursor, err := collection.Find(ctx, filter, opts)
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

func (s *UserStorage) Count(ctx context.Context) (int64, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (s *UserStorage) CountWithFilter(ctx context.Context, search string) (int64, error) {
	collection := s.appCtx.GetDatabase().MongoDB.GetCollection("users")

	// Build filter for search
	filter := bson.M{}
	if search != "" {
		// Case-insensitive search in name and email
		searchRegex := bson.M{"$regex": search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": searchRegex},
			{"email": searchRegex},
		}
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return total, nil
}
