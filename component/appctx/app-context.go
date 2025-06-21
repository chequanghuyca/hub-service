package appctx

import "hub-service/component/database"

type AppContext interface {
	SecretKey() string
	GetDatabase() *database.Database
}

type appCtx struct {
	secretKey string
	database  *database.Database
}

func NewAppContext(secretKey string, database *database.Database) AppContext {
	return &appCtx{
		secretKey: secretKey,
		database:  database,
	}
}

func (context *appCtx) SecretKey() string {
	return context.secretKey
}

func (context *appCtx) GetDatabase() *database.Database {
	return context.database
}
