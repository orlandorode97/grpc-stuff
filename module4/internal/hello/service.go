package hello

import (
	"context"
	"fmt"

	"github.com/orlandorode97/grpc-stuff/module4/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	// UnimplementedHelloServiceServer has all interfaces for all defined rpc methods.
	// In case that `Server` does not implement any rpc methods, we can take the UnimplementedHelloServiceServer
	// rpc method to keep compatibility.
	proto.UnimplementedHelloServiceServer
}

func (s *Service) SayHello(ctx context.Context, req *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	if req.GetMessage() == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	}

	return &proto.SayHelloResponse{
		Message: fmt.Sprintf("Hello: %s", req.GetMessage()),
	}, nil
}
