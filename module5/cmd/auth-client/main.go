package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/orlandorode97/grpc-stuff/module5/internal/auth"
	"github.com/orlandorode97/grpc-stuff/module5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()
	jwt, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET is reqired")
	}
	auth, err := auth.NewService(jwt)
	if err != nil {
		log.Fatal(err)
	}

	token, err := auth.IssueToken(ctx, "user-id-12345")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	md := metadata.Pairs("authorization", token)

	ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := client.Protected(ctx, &proto.ProtectedRequest{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("response received: %s\n", res.GetUserId())
}
