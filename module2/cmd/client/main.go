package main

import (
	"context"
	"fmt"
	"log"

	"github.com/orlandorode97/grpc-stuff/module2/proto"
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

	client := proto.NewHelloServiceClient(conn)

	res, err := client.SayHello(context.Background(), &proto.SayHelloRequest{})

	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			log.Fatal("unable to parse error", status.Err().Error())
		}

		log.Fatal(status.Message())
	}

	fmt.Printf("response received: %s\n", res.Message)
}
