package hello

import (
	"context"
	"errors"

	"github.com/orlandorode97/grpc-stuff/module5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const userIDCtxKey = "user_id"

type Validator interface {
	IssueToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type Middleware struct {
	validator Validator
}

func NewMiddleware(validator Validator) (*Middleware, error) {
	if validator == nil {
		return nil, errors.New("validator cannot be nil")
	}

	return &Middleware{validator: validator}, nil
}

func (m *Middleware) UnaryAuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// check RPC method when it's protected
	if info.FullMethod != proto.HelloService_Protected_FullMethodName {
		return handler(ctx, req)
	}
	// get token from metadata
	token, err := getTokenFromMetata(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token must be present")
	}
	// call validate token
	userID, err := m.validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "token is not valid: %v", err.Error())
	}
	// add user_id to the context
	ctx = context.WithValue(ctx, userIDCtxKey, userID)

	// call handler
	return handler(ctx, req)
}

func getTokenFromMetata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) != 1 {
		return "", errors.New("unable to get metdata from context")
	}

	return md["authorization"][0], nil
}
