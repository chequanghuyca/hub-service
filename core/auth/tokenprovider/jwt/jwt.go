package jwt

import (
	"hub-service/core/auth/tokenprovider"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewProvider creates a new token provider that implements the Provider interface.
func NewProvider(secret string) tokenprovider.Provider {
	return &jwtProvider{secret: secret}
}

// jwtProvider is a JWT-based token provider.
type jwtProvider struct {
	secret string
}

// NewJWTProvider creates a new jwtProvider.
func NewJWTProvider(secret string) *jwtProvider {
	return &jwtProvider{secret: secret}
}

// customClaims represents the JWT claims.
type customClaims struct {
	Payload struct {
		UserID    interface{} `json:"user_id"`
		Role      string      `json:"role"`
		Email     string      `json:"email,omitempty"`
		FirstName string      `json:"first_name,omitempty"`
		LastName  string      `json:"last_name,omitempty"`
	} `json:"payload"`
	jwt.StandardClaims
}

// Generate generates a new token pair.
func (p *jwtProvider) Generate(data tokenprovider.TokenPayload, expiry int) (*tokenprovider.Token, error) {
	// Access Token
	accessTokenClaims := &customClaims{
		Payload: struct {
			UserID    interface{} `json:"user_id"`
			Role      string      `json:"role"`
			Email     string      `json:"email,omitempty"`
			FirstName string      `json:"first_name,omitempty"`
			LastName  string      `json:"last_name,omitempty"`
		}{
			UserID:    data.UserID,
			Role:      data.Role,
			Email:     data.Email,
			FirstName: data.FirstName,
			LastName:  data.LastName,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(p.secret))
	if err != nil {
		return nil, tokenprovider.ErrEncodingToken
	}

	// Refresh Token (longer expiry, can have simpler payload if needed)
	refreshTokenClaims := &customClaims{
		Payload: struct {
			UserID    interface{} `json:"user_id"`
			Role      string      `json:"role"`
			Email     string      `json:"email,omitempty"`
			FirstName string      `json:"first_name,omitempty"`
			LastName  string      `json:"last_name,omitempty"`
		}{
			UserID:    data.UserID,
			Role:      data.Role,
			Email:     data.Email,
			FirstName: data.FirstName,
			LastName:  data.LastName,
		},
		StandardClaims: jwt.StandardClaims{
			// Refresh token expires in 2 days
			ExpiresAt: time.Now().Add(time.Hour * 24 * 2).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(p.secret))
	if err != nil {
		return nil, tokenprovider.ErrEncodingToken
	}

	return &tokenprovider.Token{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		CreatedAt:    time.Now(),
		Expiry:       expiry,
	}, nil
}

// Validate validates a token.
func (p *jwtProvider) Validate(tokenStr string) (*tokenprovider.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, tokenprovider.ErrInvalidToken
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil, tokenprovider.ErrInvalidToken
			}
		}
		return nil, tokenprovider.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		// Convert the flexible UserID back to primitive.ObjectID
		var userID primitive.ObjectID
		switch v := claims.Payload.UserID.(type) {
		case float64:
			// Handle numeric user_id (like 0)
			userID = primitive.NilObjectID
		case string:
			// Handle string user_id
			if objID, err := primitive.ObjectIDFromHex(v); err == nil {
				userID = objID
			} else {
				userID = primitive.NilObjectID
			}
		case primitive.ObjectID:
			userID = v
		default:
			userID = primitive.NilObjectID
		}

		return &tokenprovider.TokenPayload{
			UserID:    userID,
			Role:      claims.Payload.Role,
			Email:     claims.Payload.Email,
			FirstName: claims.Payload.FirstName,
			LastName:  claims.Payload.LastName,
		}, nil
	}

	return nil, tokenprovider.ErrInvalidToken
}

// GenerateAccessToken generates only a new access token.
func (p *jwtProvider) GenerateAccessToken(data tokenprovider.TokenPayload, expiry int) (string, error) {
	// Access Token
	accessTokenClaims := &customClaims{
		Payload: struct {
			UserID    interface{} `json:"user_id"`
			Role      string      `json:"role"`
			Email     string      `json:"email,omitempty"`
			FirstName string      `json:"first_name,omitempty"`
			LastName  string      `json:"last_name,omitempty"`
		}{
			UserID:    data.UserID,
			Role:      data.Role,
			Email:     data.Email,
			FirstName: data.FirstName,
			LastName:  data.LastName,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(p.secret))
	if err != nil {
		return "", tokenprovider.ErrEncodingToken
	}

	return accessTokenString, nil
}
