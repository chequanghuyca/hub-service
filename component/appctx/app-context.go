package appctx

type AppContext interface {
	SecretKey() string
}

type appCtx struct {
	secretKey string
}

func NewAppContext(secretKey string) AppContext {
	return &appCtx{
		secretKey: secretKey,
	}
}

func (context *appCtx) SecretKey() string {
	return context.secretKey
}
