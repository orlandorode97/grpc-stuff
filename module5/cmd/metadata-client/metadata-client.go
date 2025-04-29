package main

import (
	"context"
	"fmt"
	"log"

	"github.com/orlandorode97/grpc-stuff/module5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	md := metadata.Pairs("x-request-id", "12345")

	// NewOutgoingContext injects the md or metadata into the context
	ctx = metadata.NewOutgoingContext(ctx, md)

	var (
		headers  = metadata.New(map[string]string{})
		trailers = metadata.New(map[string]string{})
	)

	res, err := client.SayHello(ctx, &proto.SayHelloRequest{
		Name: "Orlando",
	}, grpc.Header(&headers),
		grpc.Trailer(&trailers),
		grpc.MaxCallRecvMsgSize(16),
		grpc.MaxCallSendMsgSize(9),
		grpc.CallContentSubtype("json"))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("response received: %s\n", res.GetMessage())
	fmt.Printf("response headers: %s\n", headers)
	fmt.Printf("response trailers: %s\n", trailers)
}
