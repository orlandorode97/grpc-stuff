package main

import (
	"log"
	"net"

	"github.com/orlandorode97/grpc-stuff/module3/exercise/internal/uploading"
	"github.com/orlandorode97/grpc-stuff/module3/exercise/proto"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	service := uploading.New()

	proto.RegisterFileUploadServiceServer(grpcServer, service)

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting grpc server on: %v\n", listen.Addr().String())

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
