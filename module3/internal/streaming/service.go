package streaming

import (
	"errors"
	"io"
	"log"
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

func (s *Service) LogStream(stream proto.StreamingService_LogStreamServer) error {
	//init count
	count := 0
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			// client closed the stream
			return stream.SendAndClose(&proto.LogStreamResponse{
				EntriesLogged: int32(count),
			})
		}

		if err != nil {
			return err
		}

		log.Printf("Received log [%s]: %s - %s\n", req.GetTimestamp().AsTime(), req.GetLevel().String(), req.GetMessage())
		count++
	}
}

func (s *Service) Echo(stream proto.StreamingService_EchoServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {

			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}
		log.Printf("message received: %s", req.GetMessage())

		res := &proto.EchoResponse{
			Message: req.GetMessage(),
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}
}
