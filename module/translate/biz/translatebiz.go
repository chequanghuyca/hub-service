package biz

import (
	"context"
	"errors"
	"fmt"
	"hub-service/core/appctx"
	challengemodel "hub-service/module/challenge/model"
	"hub-service/module/translate/model"
	"hub-service/utils/helper"
	"log"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bounoable/deepl"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChallengeStore defines the interface for challenge data access.
type ChallengeStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*challengemodel.Challenge, error)
}

// TranslateBiz handles the business logic for translation scoring.
type TranslateBiz struct {
	appCtx appctx.AppContext
	store  ChallengeStore
}

// NewTranslateBiz creates a new TranslateBiz instance.
func NewTranslateBiz(appCtx appctx.AppContext, store ChallengeStore) *TranslateBiz {
	return &TranslateBiz{
		appCtx: appCtx,
		store:  store,
	}
}

// ScoreTranslation scores a user's translation against DeepL's translation.
func (biz *TranslateBiz) ScoreTranslation(ctx context.Context, req model.ScoreRequest) (*model.ScoreResponse, error) {
	challengeID, err := primitive.ObjectIDFromHex(req.ChallengeID)
	if err != nil {
		return nil, errors.New("invalid challenge ID format")
	}

	challenge, err := biz.store.Get(ctx, challengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}

	// The challenge content can contain multiple sentences.
	// For now, we treat the whole content as one sentence and sentence_index must be 0.
	if req.SentenceIndex != 0 {
		return nil, errors.New("sentence splitting not implemented, please use sentence_index 0")
	}
	originalSentence := challenge.Content

	deeplClient := biz.appCtx.GetDeeplClient()
	if deeplClient == nil {
		return nil, errors.New("DeepL client is not configured")
	}

	// Get translation from DeepL
	deeplTranslation, _, err := deeplClient.Translate(
		ctx,
		originalSentence,
		deepl.Language(challenge.TargetLang),
		deepl.SourceLang(deepl.Language(challenge.SourceLang)),
	)
	if err != nil {
		log.Printf("DeepL API error: %v", err)
		return nil, errors.New("failed to get translation from DeepL")
	}

	// Calculate score using Sørensen-Dice coefficient
	score := strutil.Similarity(
		helper.NormalizeString(req.UserTranslation),
		helper.NormalizeString(deeplTranslation),
		metrics.NewSorensenDice(),
	)

	return &model.ScoreResponse{
		Score:            score * 100, // as percentage
		UserTranslation:  req.UserTranslation,
		DeepLTranslation: deeplTranslation,
		OriginalSentence: originalSentence,
	}, nil
}
