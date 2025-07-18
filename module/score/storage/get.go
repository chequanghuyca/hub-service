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
		{"$group": bson.M{
			"_id":              "$user_id",
			"total_score":      bson.M{"$sum": "$best_score"},
			"total_challenges": bson.M{"$sum": 1},
			"average_score":    bson.M{"$avg": "$best_score"},
			"best_score":       bson.M{"$max": "$best_score"},
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

// GetUserScoresBySection gets user scores for challenges in a specific section
func (s *Storage) GetUserScoresBySection(ctx context.Context, userID, sectionID primitive.ObjectID) ([]model.Score, error) {
	collection := s.db.MongoDB.GetCollection(model.CollectionName)

	// First, get all challenge IDs for the section
	challengeCollection := s.db.MongoDB.GetCollection("challenges")
	challengeFilter := bson.M{"section_id": sectionID}

	challengeCursor, err := challengeCollection.Find(ctx, challengeFilter)
	if err != nil {
		return nil, err
	}
	defer challengeCursor.Close(ctx)

	var challenges []bson.M
	if err = challengeCursor.All(ctx, &challenges); err != nil {
		return nil, err
	}

	if len(challenges) == 0 {
		return []model.Score{}, nil
	}

	// Extract challenge IDs
	var challengeIDs []primitive.ObjectID
	for _, challenge := range challenges {
		if id, ok := challenge["_id"].(primitive.ObjectID); ok {
			challengeIDs = append(challengeIDs, id)
		}
	}

	// Get scores for these challenges
	filter := bson.M{
		"user_id":      userID,
		"challenge_id": bson.M{"$in": challengeIDs},
	}

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

// GetUserSectionScoreSummary gets a summary of user's scores for a specific section
func (s *Storage) GetUserSectionScoreSummary(ctx context.Context, userID, sectionID primitive.ObjectID) (*model.UserScoreSummary, error) {
	scores, err := s.GetUserScoresBySection(ctx, userID, sectionID)
	if err != nil {
		return nil, err
	}

	if len(scores) == 0 {
		return &model.UserScoreSummary{
			UserID:          userID,
			TotalScore:      0,
			TotalChallenges: 0,
			AverageScore:    0,
			BestScore:       0,
		}, nil
	}

	var totalScore float64
	var bestScore float64
	totalChallenges := len(scores)

	for _, score := range scores {
		effectiveScore := score.BestScore
		// Nếu BestScore = 0 thì bỏ qua (không dùng DeepLScore nữa)
		if effectiveScore == 0 {
			continue
		}
		totalScore += effectiveScore
		if effectiveScore > bestScore {
			bestScore = effectiveScore
		}
	}

	averageScore := 0.0
	if totalChallenges > 0 {
		averageScore = totalScore / float64(totalChallenges)
	}

	return &model.UserScoreSummary{
		UserID:          userID,
		TotalScore:      totalScore,
		TotalChallenges: totalChallenges,
		AverageScore:    averageScore,
		BestScore:       bestScore,
	}, nil
}
