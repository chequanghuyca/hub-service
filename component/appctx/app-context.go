package appctx

import (
	"hub-service/component/database"
	"hub-service/component/deepl"
	"hub-service/component/tokenprovider"
	"hub-service/component/tokenprovider/jwt"

	deeplgo "github.com/bounoable/deepl"
)

type AppContext interface {
	GetSecretKey() string
	GetDatabase() *database.Database
	GetTokenProvider() tokenprovider.Provider
	GetDeeplClient() *deeplgo.Client
}

type appContext struct {
	secretKey     string
	db            *database.Database
	tokenProvider tokenprovider.Provider
	deeplClient   *deeplgo.Client
}

func NewAppContext(secretKey string, db *database.Database) *appContext {
	tokenProvider := jwt.NewProvider(secretKey)
	deeplClient := deepl.NewClient()
	return &appContext{
		secretKey:     secretKey,
		db:            db,
		tokenProvider: tokenProvider,
		deeplClient:   deeplClient,
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

func (ctx *appContext) GetDeeplClient() *deeplgo.Client {
	return ctx.deeplClient
}
