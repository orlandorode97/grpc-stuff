package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/orlandorode97/grpc-stuff/module3/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewStreamingServiceClient(conn)

	stream, err := client.LogStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for i := range 5 {
		req := &proto.LogStreamRequest{
			Timestamp: timestamppb.New(time.Now()),
			Level:     proto.LogLevel_LOG_LEVEL_INFO,
			Message:   fmt.Sprintf("Log: %d", i),
		}

		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Numbers of logs sent: %d\n", res.GetEntriesLogged())

}
