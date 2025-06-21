package appctx

import (
	"hub-service/component/database"
	"hub-service/component/tokenprovider"
	"hub-service/component/tokenprovider/jwt"
)

type AppContext interface {
	GetSecretKey() string
	GetDatabase() *database.Database
	GetTokenProvider() tokenprovider.Provider
}

type appContext struct {
	secretKey     string
	db            *database.Database
	tokenProvider tokenprovider.Provider
}

func NewAppContext(secretKey string, db *database.Database) *appContext {
	tokenProvider := jwt.NewProvider(secretKey)
	return &appContext{
		secretKey:     secretKey,
		db:            db,
		tokenProvider: tokenProvider,
	}
}

func (ctx *appContext) GetSecretKey() string {
	return ctx.secretKey
}

func (ctx *appContext) GetDatabase() *database.Database {
	return ctx.db
}

func (ctx *appContext) GetTokenProvider() tokenprovider.Provider {
	return ctx.tokenProvider
}
