package tokenprovider

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtProvider struct {
	secret string
}

func NewJWTProvider(secret string) Provider {
	return &jwtProvider{secret: secret}
}

type myClaims struct {
	Payload AccessTokenPayload `json:"payload"`
	jwt.StandardClaims
}

func (j *jwtProvider) Generate(data AccessTokenPayload, expiry int) (*Token, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		data,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	})

	myToken, err := t.SignedString([]byte(j.secret))

	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken: myToken,
		CreatedAt:   time.Now(),
		Expiry:      expiry,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (*AccessTokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !res.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return &claims.Payload, nil
}
