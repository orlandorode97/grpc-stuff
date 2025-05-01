package main

import (
	"log"
	"net"
	"os"

	"github.com/orlandorode97/grpc-stuff/module5/exercise/internal"
	"github.com/orlandorode97/grpc-stuff/module5/exercise/proto"
	"google.golang.org/grpc"
)

func main() {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		panic("JWT_SECRET is required")
	}

	t, err := internal.NewTokenService(jwtSecret)
	if err != nil {
		panic(err)
	}

	middleware := internal.NewMiddleware(t)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryMiddleware),
	)

	tokenServer := &internal.Server{}

	proto.RegisterTokenServiceServer(grpcServer, tokenServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	log.Printf("Starting grpc ser er on address: %v\n", listener.Addr().String())

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}

}
