package streaming

import (
	"time"

	"github.com/orlandorode97/grpc-stuff/module3/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	// UnimplementedHelloServiceServer has all interfaces for all defined rpc methods.
	// In case that `Server` does not implement any rpc methods, we can take the UnimplementedHelloServiceServer
	// rpc method to keep compatibility.
	proto.UnimplementedStreamingServiceServer
}

func (s *Service) StreamServerTime(req *proto.StreamServerTimeRequest, stream proto.StreamingService_StreamServerTimeServer) error {
	if req.GetIntervalSeconds() == 0 {
		return status.Error(codes.InvalidArgument, "interval must be set")
	}
	// initially ticker
	interval := time.Duration(req.GetIntervalSeconds()) * time.Second
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			currentTime := time.Now()
			err := stream.Send(&proto.StreamServerTimeResponse{
				CurrentTime: timestamppb.New(currentTime),
			})
			if err != nil {
				return err
			}
		}
	}
}
