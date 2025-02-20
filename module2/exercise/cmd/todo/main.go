package main

import (
	"log"
	"net"

	"github.com/orlandorode97/grpc-stuff/module2/exercise/internal/todo"
	"github.com/orlandorode97/grpc-stuff/module2/exercise/proto"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	todoService := todo.New()
	proto.RegisterTodoServiceServer(grpcServer, todoService)

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting grpcServer on: %v\n", listen.Addr().String())

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
