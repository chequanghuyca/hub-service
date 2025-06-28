package storage

import (
	"context"
	scoreStorage "hub-service/module/score/storage"
	"hub-service/module/section/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) Get(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*model.SectionWithChallenges, error) {
	var section model.Section

	// Get section
	collection := s.db.MongoDB.GetCollection(model.SectionName)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&section)
	if err != nil {
		return nil, err
	}

	// Get related challenges
	var challenges []model.Challenge
	challengeCollection := s.db.MongoDB.GetCollection("challenges")
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := challengeCollection.Find(ctx, bson.M{"section_id": id}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &challenges); err != nil {
		return nil, err
	}

	// Get user score for this section
	var userScore *model.UserScoreSummary
	scoreStore := scoreStorage.NewStorage(s.db)
	scoreSummary, err := scoreStore.GetUserSectionScoreSummary(ctx, userID, id)
	if err == nil && scoreSummary != nil {
		userScore = &model.UserScoreSummary{
			UserID:          scoreSummary.UserID,
			TotalScore:      scoreSummary.TotalScore,
			TotalChallenges: scoreSummary.TotalChallenges,
			AverageScore:    scoreSummary.AverageScore,
			BestScore:       scoreSummary.BestScore,
		}
	}

	result := &model.SectionWithChallenges{
		Section:    section,
		Challenges: challenges,
		UserScore:  userScore,
	}

	return result, nil
}

// GetSectionOnly returns only the section without challenges (for update operations)
func (s *Storage) GetSectionOnly(ctx context.Context, id primitive.ObjectID) (*model.Section, error) {
	var section model.Section

	collection := s.db.MongoDB.GetCollection(model.SectionName)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&section)
	if err != nil {
		return nil, err
	}

	return &section, nil
}
