package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const ClaimID = "id"

type service struct {
	secret []byte
}

func NewService(secret string) (*service, error) {
	if secret == "" {
		return nil, errors.New("secret must not be empty")
	}
	return &service{
		secret: []byte(secret),
	}, nil
}

func (s *service) IssueToken(ctx context.Context, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		ClaimID: userID,
	})

	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return signed, nil
}

func (s *service) ValidateToken(ctx context.Context, token string) (string, error) {

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return s.secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("unable to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !parsedToken.Valid || !ok {
		return "", errors.New("failed to extract claims")
	}

	id, ok := claims[ClaimID].(string)
	if !ok {
		return "", errors.New("unable get user id from claims")
	}

	return id, nil
}
