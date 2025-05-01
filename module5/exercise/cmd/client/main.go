package main

import (
	"context"
	"log"

	"github.com/orlandorode97/grpc-stuff/module5/exercise/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	var (
		ctx   = context.Background()
		token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDYwMzcxMjQsImV4cCI6MTc3NzU3MzEyNCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiT3JsYW5kbyBSb21vIEVtcGFuYWRhcyIsIm5hbWUiOiJDaHJpcyIsInJvbGUiOiJhZG1pbiJ9.bHmZv5OByNdkg2_-6d-F2J7VsrxLkHc7B3tHBsx1EPM"
	)

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewTokenServiceClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"authorization", token,
	))
	res, err := client.Validate(ctx, &proto.ValidateRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response received: %v", res.Claims)
}
