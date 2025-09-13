package biz

import (
	"context"
	"strings"
	"time"

	"hub-service/module/translation/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTranslationStore interface {
	CreateTranslation(ctx context.Context, data *model.TranslationCreate) error
	CreateSentence(ctx context.Context, data *model.TranslationSentenceCreate) error
}

type createTranslationBiz struct {
	store   CreateTranslationStore
	apiKey  string
	baseURL string
}

func NewCreateTranslationBiz(store CreateTranslationStore, apiKey, baseURL string) *createTranslationBiz {
	return &createTranslationBiz{
		store:   store,
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (biz *createTranslationBiz) CreateTranslation(ctx context.Context, data *model.TranslationCreate) (*model.Translation, error) {
	// Split content into sentences using improved splitter
	splitter := NewSentenceSplitter()
	sentences := splitter.SplitIntoSentencesAdvanced(data.Content)

	// Calculate total score (each sentence worth 10 points by default)
	totalScore := float64(len(sentences)) * 10.0

	// Create translation
	translation := &model.Translation{
		ID:         primitive.NewObjectID(),
		Title:      data.Title,
		Content:    data.Content,
		SourceLang: data.SourceLang,
		TargetLang: data.TargetLang,
		Category:   data.Category,
		Difficulty: data.Difficulty,
		TotalScore: totalScore,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Image:      data.Image,
	}

	// Create translation record
	translationCreate := &model.TranslationCreate{
		ID:         translation.ID,
		Title:      translation.Title,
		Content:    translation.Content,
		SourceLang: translation.SourceLang,
		TargetLang: translation.TargetLang,
		Category:   translation.Category,
		Difficulty: translation.Difficulty,
		CreatedAt:  &translation.CreatedAt,
		UpdatedAt:  &translation.UpdatedAt,
		Image:      translation.Image,
	}

	if err := biz.store.CreateTranslation(ctx, translationCreate); err != nil {
		return nil, err
	}

	// Create sentences
	for i, sentence := range sentences {
		sentenceData := &model.TranslationSentenceCreate{
			ID:            primitive.NewObjectID(),
			TranslationID: translation.ID,
			SentenceIndex: i,
			Content:       strings.TrimSpace(sentence),
			MaxScore:      10.0, // Each sentence worth 10 points
			CreatedAt:     &translation.CreatedAt,
			UpdatedAt:     &translation.UpdatedAt,
		}

		if err := biz.store.CreateSentence(ctx, sentenceData); err != nil {
			return nil, err
		}
	}

	return translation, nil
}
