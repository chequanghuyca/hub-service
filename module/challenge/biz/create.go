package biz

import (
	"context"
	"hub-service/module/challenge/model"
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
	if err := biz.store.Create(ctx, data); err != nil {
		return err
	}
	return nil
}
