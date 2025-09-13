package storage

import (
	"context"
	"time"

	common "hub-service/common"
	translationmodel "hub-service/module/translation/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Translation operations
func (s *Storage) CreateTranslation(ctx context.Context, data *translationmodel.TranslationCreate) error {
	now := time.Now()
	data.CreatedAt = &now
	data.UpdatedAt = &now

	collection := s.db.MongoDB.Database.Collection(translationmodel.TranslationCollectionName)
	_, err := collection.InsertOne(ctx, data)
	return err
}

func (s *Storage) GetTranslation(ctx context.Context, id primitive.ObjectID) (*translationmodel.Translation, error) {
	collection := s.db.MongoDB.Database.Collection(translationmodel.TranslationCollectionName)

	var translation translationmodel.Translation
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&translation)
	if err != nil {
		return nil, err
	}

	return &translation, nil
}

func (s *Storage) ListTranslations(ctx context.Context, filter map[string]interface{}, paging *common.Paging) ([]translationmodel.Translation, error) {
	collection := s.db.MongoDB.Database.Collection(translationmodel.TranslationCollectionName)

	opts := options.Find()
	if paging != nil {
		opts.SetLimit(int64(paging.Limit))
		opts.SetSkip(int64((paging.Page - 1) * paging.Limit))
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var translations []translationmodel.Translation
	if err = cursor.All(ctx, &translations); err != nil {
		return nil, err
	}

	return translations, nil
}

func (s *Storage) UpdateTranslation(ctx context.Context, id primitive.ObjectID, data *translationmodel.TranslationUpdate) error {
	now := time.Now()
	data.UpdatedAt = &now

	collection := s.db.MongoDB.Database.Collection(translationmodel.TranslationCollectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	return err
}

func (s *Storage) DeleteTranslation(ctx context.Context, id primitive.ObjectID) error {
	collection := s.db.MongoDB.Database.Collection(translationmodel.TranslationCollectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Sentence operations
func (s *Storage) CreateSentence(ctx context.Context, data *translationmodel.TranslationSentenceCreate) error {
	now := time.Now()
	data.CreatedAt = &now
	data.UpdatedAt = &now

	collection := s.db.MongoDB.Database.Collection(translationmodel.SentenceCollectionName)
	_, err := collection.InsertOne(ctx, data)
	return err
}

func (s *Storage) GetSentencesByTranslationID(ctx context.Context, translationID primitive.ObjectID) ([]translationmodel.TranslationSentence, error) {
	collection := s.db.MongoDB.Database.Collection(translationmodel.SentenceCollectionName)

	opts := options.Find().SetSort(bson.D{{Key: "sentence_index", Value: 1}})
	cursor, err := collection.Find(ctx, bson.M{"translation_id": translationID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sentences []translationmodel.TranslationSentence
	if err = cursor.All(ctx, &sentences); err != nil {
		return nil, err
	}

	return sentences, nil
}

func (s *Storage) UpdateSentence(ctx context.Context, id primitive.ObjectID, data *translationmodel.TranslationSentenceCreate) error {
	now := time.Now()
	data.UpdatedAt = &now

	collection := s.db.MongoDB.Database.Collection(translationmodel.SentenceCollectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	return err
}

func (s *Storage) DeleteSentence(ctx context.Context, id primitive.ObjectID) error {
	collection := s.db.MongoDB.Database.Collection(translationmodel.SentenceCollectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *Storage) DeleteSentencesByTranslationID(ctx context.Context, translationID primitive.ObjectID) error {
	collection := s.db.MongoDB.Database.Collection(translationmodel.SentenceCollectionName)
	_, err := collection.DeleteMany(ctx, bson.M{"translation_id": translationID})
	return err
}

// User score operations
func (s *Storage) CreateUserScore(ctx context.Context, data *translationmodel.UserTranslationScoreCreate) error {
	collection := s.db.MongoDB.Database.Collection(translationmodel.UserTranslationScoreCollectionName)
	_, err := collection.InsertOne(ctx, data)
	return err
}

func (s *Storage) GetUserScore(ctx context.Context, userID, translationID primitive.ObjectID, sentenceIndex int) (*translationmodel.UserTranslationScore, error) {
	collection := s.db.MongoDB.Database.Collection(translationmodel.UserTranslationScoreCollectionName)

	var score translationmodel.UserTranslationScore
	err := collection.FindOne(ctx, bson.M{
		"user_id":        userID,
		"translation_id": translationID,
		"sentence_index": sentenceIndex,
	}).Decode(&score)
	if err != nil {
		return nil, err
	}

	return &score, nil
}

func (s *Storage) UpdateUserScore(ctx context.Context, id primitive.ObjectID, data *translationmodel.UserTranslationScoreCreate) error {
	collection := s.db.MongoDB.Database.Collection(translationmodel.UserTranslationScoreCollectionName)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	return err
}

func (s *Storage) GetUserScoresByTranslation(ctx context.Context, userID, translationID primitive.ObjectID) ([]translationmodel.UserTranslationScore, error) {
	collection := s.db.MongoDB.Database.Collection(translationmodel.UserTranslationScoreCollectionName)

	opts := options.Find().SetSort(bson.D{{Key: "sentence_index", Value: 1}})
	cursor, err := collection.Find(ctx, bson.M{
		"user_id":        userID,
		"translation_id": translationID,
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scores []translationmodel.UserTranslationScore
	if err = cursor.All(ctx, &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

func (s *Storage) GetUserTranslationSummaries(ctx context.Context, userID primitive.ObjectID) ([]translationmodel.TranslationSummary, error) {
	collection := s.db.MongoDB.Database.Collection(translationmodel.UserTranslationScoreCollectionName)

	// Aggregate pipeline to get summaries
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":              "$translation_id",
			"total_sentences":  bson.M{"$sum": 1},
			"total_user_score": bson.M{"$sum": "$best_score"},
			"last_attempt_at":  bson.M{"$max": "$updated_at"},
		}},
		{"$lookup": bson.M{
			"from":         translationmodel.TranslationCollectionName,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "translation",
		}},
		{"$unwind": "$translation"},
		{"$project": bson.M{
			"translation_id":     "$_id",
			"title":              "$translation.title",
			"total_sentences":    "$total_sentences",
			"total_user_score":   "$total_user_score",
			"max_possible_score": "$translation.total_score",
			"last_attempt_at":    "$last_attempt_at",
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []translationmodel.TranslationSummary
	if err = cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}

	// Calculate progress percent
	for i := range summaries {
		if summaries[i].MaxPossibleScore > 0 {
			summaries[i].ProgressPercent = (summaries[i].TotalUserScore / summaries[i].MaxPossibleScore) * 100
		}
		summaries[i].CompletedCount = int(summaries[i].TotalSentences)
	}

	return summaries, nil
}
