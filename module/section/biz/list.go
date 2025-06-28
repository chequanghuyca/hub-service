package biz

import (
	"context"
	"hub-service/common"
	"hub-service/module/section/model"
)

type ListSectionStore interface {
	List(ctx context.Context, paging *common.Paging, moreKeys ...string) ([]model.Section, error)
}

type listSectionBiz struct {
	store ListSectionStore
}

func NewListSectionBiz(store ListSectionStore) *listSectionBiz {
	return &listSectionBiz{store: store}
}

func (biz *listSectionBiz) ListSection(ctx context.Context, paging *common.Paging) ([]model.Section, error) {
	return biz.store.List(ctx, paging)
}
