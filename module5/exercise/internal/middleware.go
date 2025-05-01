package internal

import (
	"context"
	"errors"

	"github.com/orlandorode97/grpc-stuff/module5/exercise/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const claimsContextKey = "claims"

type Middleware struct {
	t *Token
}

func NewMiddleware(token *Token) *Middleware {
	return &Middleware{
		t: token,
	}
}

func (m *Middleware) UnaryMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if info.FullMethod != proto.TokenService_Validate_FullMethodName {
		return handler(ctx, req)
	}

	token, err := getTokenMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token must be present")
	}

	claims, err := m.t.Validate(ctx, token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	ctx = context.WithValue(ctx, claimsContextKey, claims)

	return handler(ctx, req)
}

func getTokenMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok || len(md["authorization"]) != 1 {
		return "", errors.New("unable to get metadata from context")
	}

	return md["authorization"][0], nil
}
