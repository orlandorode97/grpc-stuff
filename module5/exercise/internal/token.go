package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Name string `json:"name"`
	Role string `json:"role"`
}

func (c Claims) Valid() error {
	return nil
}

type Token struct {
	secret []byte
}

func NewTokenService(secret string) (*Token, error) {
	if secret == "" {
		return nil, errors.New("secret is required")
	}
	return &Token{
		secret: []byte(secret),
	}, nil
}

func (t *Token) Validate(ctx context.Context, token string) (*Claims, error) {
	c := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, c, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}

		return t.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}

	return c, nil
}
