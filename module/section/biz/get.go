package biz

import (
	"context"
	"hub-service/module/section/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetSectionStore interface {
	Get(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*model.SectionWithChallenges, error)
}

type getSectionBiz struct {
	store GetSectionStore
}

func NewGetSectionBiz(store GetSectionStore) *getSectionBiz {
	return &getSectionBiz{store: store}
}

func (biz *getSectionBiz) GetSection(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*model.SectionWithChallenges, error) {
	return biz.store.Get(ctx, id, userID)
}
