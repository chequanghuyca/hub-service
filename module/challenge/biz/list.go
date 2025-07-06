package biz

import (
	"context"
	"hub-service/common"
	"hub-service/module/challenge/model"
)

type ListChallengeStore interface {
	List(
		ctx context.Context,
		paging *common.Paging,
		sectionID string,
		search string,
		moreKeys ...string) ([]model.Challenge, error)
}

type listChallengeBiz struct {
	store ListChallengeStore
}

func NewListChallengeBiz(store ListChallengeStore) *listChallengeBiz {
	return &listChallengeBiz{store: store}
}

func (biz *listChallengeBiz) ListChallenge(ctx context.Context, paging *common.Paging, sectionID string, search string) ([]model.Challenge, error) {
	result, err := biz.store.List(ctx, paging, sectionID, search)
	if err != nil {
		return nil, err
	}
	return result, nil
}
