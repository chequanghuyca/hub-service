package biz

import (
	"context"
	"hub-service/module/challenge/model"
	"slices"
)

type CreateChallengeStore interface {
	Create(ctx context.Context, data *model.ChallengeCreate) error
}

type createChallengeBiz struct {
	store CreateChallengeStore
}

func NewCreateChallengeBiz(store CreateChallengeStore) *createChallengeBiz {
	return &createChallengeBiz{store: store}
}

func (biz *createChallengeBiz) CreateChallenge(ctx context.Context, data *model.ChallengeCreate) error {
	// Additional business validation using constants
	if err := biz.validateDifficulty(data.Difficulty); err != nil {
		return err
	}

	if err := biz.validateCategory(data.Category); err != nil {
		return err
	}

	if err := biz.store.Create(ctx, data); err != nil {
		return err
	}
	return nil
}

// validateDifficulty validates difficulty using constants
func (biz *createChallengeBiz) validateDifficulty(difficulty string) error {
	if !slices.Contains(model.GetValidDifficulties(), difficulty) {
		return model.ErrInvalidDifficulty
	}
	return nil
}

// validateCategory validates category using constants
func (biz *createChallengeBiz) validateCategory(category string) error {
	if !slices.Contains(model.GetValidCategories(), category) {
		return model.ErrInvalidCategory
	}
	return nil
}
