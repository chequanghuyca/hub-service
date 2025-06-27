package storage

import (
	"context"
	"hub-service/module/score/model"
	"hub-service/utils/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) GetScoreByUserAndChallenge(ctx context.Context, userID, challengeID primitive.ObjectID) (*model.Score, error) {
	collection := s.db.MongoDB.GetCollection(model.CollectionName)

	filter := bson.M{
		"user_id":      userID,
		"challenge_id": challengeID,
	}

	var score model.Score
	err := collection.FindOne(ctx, filter).Decode(&score)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &score, nil
}

func (s *Storage) GetUserScores(ctx context.Context, userID primitive.ObjectID) ([]model.Score, error) {
	collection := s.db.MongoDB.GetCollection(model.CollectionName)

	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []model.Score
	if err = cursor.All(ctx, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

func (s *Storage) GetUserScoreSummary(ctx context.Context, userID primitive.ObjectID) (*model.UserScoreSummary, error) {
	collection := s.db.MongoDB.GetCollection(model.CollectionName)

	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$addFields": bson.M{
			"effective_score": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$gt": []interface{}{"$best_score", 0}},
					"then": "$best_score",
					"else": "$deepl_score",
				},
			},
		}},
		{"$group": bson.M{
			"_id":              "$user_id",
			"total_score":      bson.M{"$sum": "$effective_score"},
			"total_challenges": bson.M{"$sum": 1},
			"average_score":    bson.M{"$avg": "$effective_score"},
			"best_score":       bson.M{"$max": "$effective_score"},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &model.UserScoreSummary{
			UserID:          userID,
			TotalScore:      0,
			TotalChallenges: 0,
			AverageScore:    0,
			BestScore:       0,
		}, nil
	}

	result := results[0]

	totalScore, _ := helper.ToFloat64(result["total_score"])
	avgScore, _ := helper.ToFloat64(result["average_score"])
	bestScore, _ := helper.ToFloat64(result["best_score"])

	return &model.UserScoreSummary{
		UserID:          userID,
		TotalScore:      totalScore,
		TotalChallenges: int(result["total_challenges"].(int32)),
		AverageScore:    avgScore,
		BestScore:       bestScore,
	}, nil
}
