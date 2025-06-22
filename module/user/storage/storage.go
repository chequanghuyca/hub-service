package storage

import (
	"hub-service/core/appctx"
)

type UserStorage struct {
	appCtx appctx.AppContext
}

func NewUserStorage(appCtx appctx.AppContext) *UserStorage {
	return &UserStorage{appCtx: appCtx}
}
