package tokenprovider

import (
	"errors"
	"hub-service/common"
	"time"
)

type Provider interface {
	Generate(data AccessTokenPayload, expiry int) (*Token, error)
	Validate(token string) (*AccessTokenPayload, error)
}

var (
	ErrNotFound = common.NewCustomError(
		errors.New("token not found"),
		"token not found",
		"ErrNotFound",
	)

	ErrEncodingToken = common.NewCustomError(
		errors.New("error encoding token"),
		"error encoding token",
		"ErrEncodingToken",
	)

	ErrInvalidToken = common.NewCustomError(
		errors.New("invalid token"),
		"invalid token",
		"ErrInvalidToken",
	)
)

type Token struct {
	AccessToken string    `json:"access_token"`
	FreshToken  string    `json:"fresh_token"`
	CreatedAt   time.Time `json:"created"`
	Expiry      int       `json:"expiry"`
}

type AccessTokenPayload struct {
	UserId    int    `json:"user_id"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}

type FreshTokenPayload struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
}
