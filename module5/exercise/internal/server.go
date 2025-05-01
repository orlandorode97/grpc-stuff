package internal

import (
	"context"

	"github.com/orlandorode97/grpc-stuff/module5/exercise/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedTokenServiceServer
}

func (s *Server) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get claism from context")
	}

	return &proto.ValidateResponse{
		Claims: &proto.Claims{
			Exp:  claims.ExpiresAt.String(),
			Ait:  claims.IssuedAt.String(),
			Name: claims.Name,
			Role: claims.Role,
			Sub:  claims.Subject,
		},
	}, nil
}
