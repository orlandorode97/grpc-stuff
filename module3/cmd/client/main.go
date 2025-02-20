package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/orlandorode97/grpc-stuff/module3/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewStreamingServiceClient(conn)

	stream, err := client.StreamServerTime(context.Background(), &proto.StreamServerTimeRequest{
		IntervalSeconds: 2,
	})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			log.Fatal("unable to parse error", status.Err().Error())
		}

		log.Fatal(status.Message())
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println(err.Error())
			break
		}
		if err != nil {
			log.Fatal(err)
			break
		}
		fmt.Println("Current server time", res.CurrentTime.AsTime())
	}
}
