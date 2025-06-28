package biz

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteSectionStore interface {
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type deleteSectionBiz struct {
	store DeleteSectionStore
}

func NewDeleteSectionBiz(store DeleteSectionStore) *deleteSectionBiz {
	return &deleteSectionBiz{store: store}
}

func (biz *deleteSectionBiz) DeleteSection(ctx context.Context, id primitive.ObjectID) error {
	return biz.store.Delete(ctx, id)
}
