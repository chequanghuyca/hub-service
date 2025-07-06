package biz

import (
	"context"
	"hub-service/common"
	"hub-service/module/section/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListSectionStore interface {
	List(ctx context.Context, paging *common.Paging, userID primitive.ObjectID, title string, moreKeys ...string) ([]model.SectionWithScore, error)
}

type listSectionBiz struct {
	store ListSectionStore
}

func NewListSectionBiz(store ListSectionStore) *listSectionBiz {
	return &listSectionBiz{store: store}
}

func (biz *listSectionBiz) ListSection(ctx context.Context, paging *common.Paging, userID primitive.ObjectID, title string) ([]model.SectionWithScore, error) {
	return biz.store.List(ctx, paging, userID, title)
}

type ListSimpleSectionStore interface {
	ListSimple(ctx context.Context, title string) ([]model.SectionSimple, error)
}

type listSimpleSectionBiz struct {
	store ListSimpleSectionStore
}

func NewListSimpleSectionBiz(store ListSimpleSectionStore) *listSimpleSectionBiz {
	return &listSimpleSectionBiz{store: store}
}

func (biz *listSimpleSectionBiz) ListSimpleSection(ctx context.Context, title string) ([]model.SectionSimple, error) {
	return biz.store.ListSimple(ctx, title)
}
