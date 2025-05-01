package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/orlandorode97/grpc-stuff/module6/internal/hello"
	"github.com/orlandorode97/grpc-stuff/module6/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	ctx := context.Background()

	serveConfig := &hello.Config{
		MethodConfig: []*hello.MethodConfig{{
			Names: []*hello.NameConfig{{
				Service: "hello.HelloService",
				Method:  "LongRunning",
			}},
			Timeout: "4s",
		}},
	}

	svcConfig, err := json.Marshal(serveConfig)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(string(svcConfig)),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	_, err = client.LongRunning(ctx, &proto.LongRunningRequest{})

	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			log.Fatal("unable to parse error", status.Err().Error())
		}

		log.Fatal(status.Message())
	}

	log.Println("successfully called RPC method")
}
