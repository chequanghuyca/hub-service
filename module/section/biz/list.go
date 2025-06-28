package biz

import (
	"context"
	"hub-service/common"
	"hub-service/module/section/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListSectionStore interface {
	List(ctx context.Context, paging *common.Paging, userID primitive.ObjectID, moreKeys ...string) ([]model.SectionWithScore, error)
}

type listSectionBiz struct {
	store ListSectionStore
}

func NewListSectionBiz(store ListSectionStore) *listSectionBiz {
	return &listSectionBiz{store: store}
}

func (biz *listSectionBiz) ListSection(ctx context.Context, paging *common.Paging, userID primitive.ObjectID) ([]model.SectionWithScore, error) {
	return biz.store.List(ctx, paging, userID)
}
