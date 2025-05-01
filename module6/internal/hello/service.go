package hello

import (
	"context"
	"log"
	"time"

	"github.com/orlandorode97/grpc-stuff/module6/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	name string
	// UnimplementedHelloServiceServer has all interfaces for all defined rpc methods.
	// In case that `Server` does not implement any rpc methods, we can take the UnimplementedHelloServiceServer
	// rpc method to keep compatibility.
	proto.UnimplementedHelloServiceServer
}

func NewService(name string) *Service {
	return &Service{
		name: name,
	}
}

func (s *Service) LongRunning(ctx context.Context, req *proto.LongRunningRequest) (*proto.LongRunningResponse, error) {
	select {
	case <-time.Tick(time.Second * 5):
		log.Printf("finished waiting, not end request successfully")
	case <-ctx.Done():
		log.Printf("context cancelled")
		return nil, ctx.Err()
	}

	return &proto.LongRunningResponse{}, nil
}

func (s *Service) Flaky(ctx context.Context, req *proto.FlakyRequest) (*proto.FlakyResponse, error) {
	log.Println("error response returned")
	return nil, status.Error(codes.Internal, "flaky error occurred")
}

func (s *Service) GetServerAddress(ctx context.Context, req *proto.GetServerAddressRequest) (*proto.GetServerAddressResponse, error) {
	log.Printf("requested received: %s", s.name)
	return &proto.GetServerAddressResponse{
		Address: s.name,
	}, nil
}
