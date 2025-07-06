package storage

import (
	"context"
	"hub-service/common"
	scoreStorage "hub-service/module/score/storage"
	"hub-service/module/section/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) List(ctx context.Context, paging *common.Paging, userID primitive.ObjectID, title string, moreKeys ...string) ([]model.SectionWithScore, error) {
	var sections []model.Section
	collection := s.db.MongoDB.GetCollection(model.SectionName)

	// Build filter query
	filter := bson.M{}

	// Add title search if provided
	if title != "" {
		filter["title"] = bson.M{
			"$regex":   title,
			"$options": "i", // Case-insensitive search
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

	if err = cursor.All(ctx, &sections); err != nil {
		return nil, err
	}

	// Get user scores for each section
	scoreStore := scoreStorage.NewStorage(s.db)
	var result []model.SectionWithScore

	for _, section := range sections {
		userScore, err := scoreStore.GetUserSectionScoreSummary(ctx, userID, section.ID)
		if err != nil {
			// If error getting score, continue without score
			result = append(result, model.SectionWithScore{
				Section:   section,
				UserScore: nil,
			})
			continue
		}

		// Convert score model to section model
		sectionUserScore := &model.UserScoreSummary{
			UserID:          userScore.UserID,
			TotalScore:      userScore.TotalScore,
			TotalChallenges: userScore.TotalChallenges,
			AverageScore:    userScore.AverageScore,
			BestScore:       userScore.BestScore,
		}

		result = append(result, model.SectionWithScore{
			Section:   section,
			UserScore: sectionUserScore,
		})
	}

	// Total count for paging with same filter
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	paging.Total = total

	return result, nil
}

// ListSimple returns all sections with only id and title, optionally filtered by title
func (s *Storage) ListSimple(ctx context.Context, title string) ([]model.SectionSimple, error) {
	var result []model.SectionSimple
	collection := s.db.MongoDB.GetCollection(model.SectionName)

	// Build filter query
	filter := bson.M{}

	// Add title search if provided
	if title != "" {
		filter["title"] = bson.M{
			"$regex":   title,
			"$options": "i", // Case-insensitive search
		}
	}

	// Only select id and title fields
	projection := bson.M{
		"_id":   1,
		"title": 1,
	}

	// Sorting by created_at desc
	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
