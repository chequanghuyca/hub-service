package biz

import (
	"context"
	"hub-service/module/section/model"
)

type CreateSectionStore interface {
	Create(ctx context.Context, data *model.SectionCreate) error
}

type createSectionBiz struct {
	store CreateSectionStore
}

func NewCreateSectionBiz(store CreateSectionStore) *createSectionBiz {
	return &createSectionBiz{store: store}
}

func (biz *createSectionBiz) CreateSection(ctx context.Context, data *model.SectionCreate) error {
	return biz.store.Create(ctx, data)
}
