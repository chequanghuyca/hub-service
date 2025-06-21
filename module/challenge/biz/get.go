package biz

import (
	"context"
	"hub-service/module/challenge/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetChallengeStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*model.Challenge, error)
}

type getChallengeBiz struct {
	store GetChallengeStore
}

func NewGetChallengeBiz(store GetChallengeStore) *getChallengeBiz {
	return &getChallengeBiz{store: store}
}

func (biz *getChallengeBiz) GetChallenge(ctx context.Context, id primitive.ObjectID) (*model.Challenge, error) {
	data, err := biz.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return data, nil
}
