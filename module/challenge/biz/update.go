package biz

import (
	"context"
	"hub-service/module/challenge/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateChallengeStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*model.Challenge, error)
	Update(ctx context.Context, id primitive.ObjectID, data *model.ChallengeUpdate) error
}

type updateChallengeBiz struct {
	store UpdateChallengeStore
}

func NewUpdateChallengeBiz(store UpdateChallengeStore) *updateChallengeBiz {
	return &updateChallengeBiz{store: store}
}

func (biz *updateChallengeBiz) UpdateChallenge(
	ctx context.Context,
	id primitive.ObjectID,
	data *model.ChallengeUpdate,
) error {
	_, err := biz.store.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := biz.store.Update(ctx, id, data); err != nil {
		return err
	}
	return nil
}
