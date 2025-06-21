package biz

import (
	"context"
	"hub-service/module/challenge/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteChallengeStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*model.Challenge, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type deleteChallengeBiz struct {
	store DeleteChallengeStore
}

func NewDeleteChallengeBiz(store DeleteChallengeStore) *deleteChallengeBiz {
	return &deleteChallengeBiz{store: store}
}

func (biz *deleteChallengeBiz) DeleteChallenge(ctx context.Context, id primitive.ObjectID) error {
	_, err := biz.store.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := biz.store.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
