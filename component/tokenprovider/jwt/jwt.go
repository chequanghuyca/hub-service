package jwt

import (
	"hub-service/component/tokenprovider"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	Payload tokenprovider.TokenPayload `json:"payload"`
	jwt.StandardClaims
}

// Generate generates a new token pair.
func (p *jwtProvider) Generate(data tokenprovider.TokenPayload, expiry int) (*tokenprovider.Token, error) {
	// Access Token
	accessTokenClaims := &customClaims{
		Payload: data,
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
		Payload: data, // Using the same payload for simplicity
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
		return &claims.Payload, nil
	}

	return nil, tokenprovider.ErrInvalidToken
}
