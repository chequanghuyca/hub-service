package biz

import (
	"context"

	"hub-service/module/translation/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetTranslationStore interface {
	GetTranslation(ctx context.Context, id primitive.ObjectID) (*model.Translation, error)
	GetSentencesByTranslationID(ctx context.Context, translationID primitive.ObjectID) ([]model.TranslationSentence, error)
	GetUserScoresByTranslation(ctx context.Context, userID, translationID primitive.ObjectID) ([]model.UserTranslationScore, error)
	GetUserTranslationSummaries(ctx context.Context, userID primitive.ObjectID) ([]model.TranslationSummary, error)
}

type getTranslationBiz struct {
	store GetTranslationStore
}

func NewGetTranslationBiz(store GetTranslationStore) *getTranslationBiz {
	return &getTranslationBiz{store: store}
}

func (biz *getTranslationBiz) GetTranslation(ctx context.Context, id primitive.ObjectID) (*model.TranslationWithSentences, error) {
	translation, err := biz.store.GetTranslation(ctx, id)
	if err != nil {
		return nil, err
	}

	sentences, err := biz.store.GetSentencesByTranslationID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.TranslationWithSentences{
		Translation: *translation,
		Sentences:   sentences,
	}, nil
}

func (biz *getTranslationBiz) GetTranslationWithUserProgress(ctx context.Context, translationID, userID primitive.ObjectID) (*model.TranslationWithUserProgress, error) {
	translation, err := biz.store.GetTranslation(ctx, translationID)
	if err != nil {
		return nil, err
	}

	sentences, err := biz.store.GetSentencesByTranslationID(ctx, translationID)
	if err != nil {
		return nil, err
	}

	userScores, err := biz.store.GetUserScoresByTranslation(ctx, userID, translationID)
	if err != nil {
		return nil, err
	}

	// Calculate progress
	totalUserScore := 0.0
	completedCount := 0
	totalPossibleScore := 0.0

	// Calculate total possible score from sentences
	for _, sentence := range sentences {
		totalPossibleScore += sentence.MaxScore
	}

	// Calculate user's total score and completed count
	for _, score := range userScores {
		totalUserScore += score.BestScore
		completedCount++
	}

	progressPercent := 0.0
	if totalPossibleScore > 0 {
		progressPercent = (totalUserScore / totalPossibleScore) * 100
	}

	return &model.TranslationWithUserProgress{
		Translation:     *translation,
		Sentences:       sentences,
		UserScores:      userScores,
		TotalUserScore:  totalUserScore,
		CompletedCount:  completedCount,
		ProgressPercent: progressPercent,
	}, nil
}

func (biz *getTranslationBiz) GetUserTranslationScores(ctx context.Context, userID primitive.ObjectID) ([]model.TranslationSummary, error) {
	return biz.store.GetUserTranslationSummaries(ctx, userID)
}
