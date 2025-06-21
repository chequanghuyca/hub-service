package tokenprovider

import (
	"errors"
	"hub-service/common"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Provider defines the interface for token generation and validation.
type Provider interface {
	// Generate generates a new token pair (access and refresh).
	Generate(data TokenPayload, expiry int) (*Token, error)
	// Validate validates a token string and returns its payload.
	Validate(token string) (*TokenPayload, error)
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

// Token represents a pair of access and refresh tokens.
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	Expiry       int       `json:"expiry"`
}

// TokenPayload contains the data stored in a token.
type TokenPayload struct {
	UserID primitive.ObjectID `json:"user_id"`
	Role   string             `json:"role"`
}
