package hello

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/orlandorode97/grpc-stuff/module5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service struct {
	// UnimplementedHelloServiceServer has all interfaces for all defined rpc methods.
	// In case that `Server` does not implement any rpc methods, we can take the UnimplementedHelloServiceServer
	// rpc method to keep compatibility.
	proto.UnimplementedHelloServiceServer
}

func (s *Service) SayHello(ctx context.Context, req *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	start := time.Now()
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md)

	if !ok || len(md["x-request-id"]) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot have metadata")
	}

	requestID := md["x-request-id"][0]
	log.Printf("request id %s\n", requestID)

	header := metadata.New(map[string]string{"request-start-timestamp": start.String()})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, status.Error(codes.Internal, "failed to send header")
	}

	trailer := metadata.New(map[string]string{"request-end-timestamp": time.Now().String()})
	if err := grpc.SetTrailer(ctx, trailer); err != nil {
		return nil, status.Error(codes.Internal, "failed to send trailer")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	}

	return &proto.SayHelloResponse{
		Message: fmt.Sprintf("Hello: %s", req.GetName()),
	}, nil
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

func (s *Service) Protected(ctx context.Context, req *proto.ProtectedRequest) (*proto.ProtectedResponse, error) {
	userID, ok := ctx.Value(userIDCtxKey).(string)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "user id missing from context")
	}

	return &proto.ProtectedResponse{
		UserId: userID,
	}, nil
}
