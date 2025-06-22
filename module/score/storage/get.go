package storage

import (
	"context"
	"fmt"
	"hub-service/module/score/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// toFloat64 safely converts an interface{} to float64, handling multiple numeric types.
func toFloat64(v interface{}) (float64, error) {
	switch i := v.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

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

	totalScore, _ := toFloat64(result["total_score"])
	avgScore, _ := toFloat64(result["average_score"])
	bestScore, _ := toFloat64(result["best_score"])

	return &model.UserScoreSummary{
		UserID:          userID,
		TotalScore:      totalScore,
		TotalChallenges: int(result["total_challenges"].(int32)),
		AverageScore:    avgScore,
		BestScore:       bestScore,
	}, nil
}

func (s *Storage) GetTotalScore(ctx context.Context, userID primitive.ObjectID) (*model.GetTotalScoreResponse, error) {
	summary, err := s.GetUserScoreSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.GetTotalScoreResponse{
		UserID:          userID.Hex(),
		TotalScore:      summary.TotalScore,
		TotalChallenges: summary.TotalChallenges,
		AverageScore:    summary.AverageScore,
		BestScore:       summary.BestScore,
	}, nil
}
