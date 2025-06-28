package biz

import (
	"context"
	"hub-service/module/section/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateSectionStore interface {
	GetSectionOnly(ctx context.Context, id primitive.ObjectID) (*model.Section, error)
	Update(ctx context.Context, id primitive.ObjectID, data *model.SectionUpdate) error
}

type updateSectionBiz struct {
	store UpdateSectionStore
}

func NewUpdateSectionBiz(store UpdateSectionStore) *updateSectionBiz {
	return &updateSectionBiz{store: store}
}

func (biz *updateSectionBiz) UpdateSection(ctx context.Context, id primitive.ObjectID, data *model.SectionUpdate) error {
	_, err := biz.store.GetSectionOnly(ctx, id)
	if err != nil {
		return err
	}

	if err := biz.store.Update(ctx, id, data); err != nil {
		return err
	}
	return nil
}
