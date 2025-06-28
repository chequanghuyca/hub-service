package storage

import (
	"context"
	"hub-service/common"
	"hub-service/module/challenge/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) List(
	ctx context.Context,
	paging *common.Paging,
	moreKeys ...string) ([]model.Challenge, error) {
	var result []model.Challenge
	collection := s.db.MongoDB.GetCollection(model.CollectionName)

	// Build filter
	filter := bson.M{}

	// Add SectionID filter if provided
	if len(moreKeys) > 0 && moreKeys[0] != "" {
		sectionID, err := primitive.ObjectIDFromHex(moreKeys[0])
		if err == nil {
			filter["section_id"] = sectionID
		}
	}

	// Paging
	findOptions := options.Find()
	findOptions.SetSkip(int64((paging.Page - 1) * paging.Limit))
	findOptions.SetLimit(int64(paging.Limit))

	// Sorting
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	// Total count for paging
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	paging.Total = total

	return result, nil
}
