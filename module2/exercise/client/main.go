package main

import (
	"context"
	"log"

	"github.com/orlandorode97/grpc-stuff/module2/exercise/proto"
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

	client := proto.NewTodoServiceClient(conn)

	_, err = client.AddTask(context.TODO(), &proto.AddTaskRequest{
		Task: "This is my task from a client",
	})

	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			log.Fatalf("unable to parse grpc status error: %w\n", err.Error())
		}
		log.Fatal(status.String())
	}

}
