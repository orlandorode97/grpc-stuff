package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/orlandorode97/grpc-stuff/module6/internal/resolve"
	"github.com/orlandorode97/grpc-stuff/module6/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var builder resolve.Builder

func init() {
	builder := resolve.Builder{}
	resolver.Register(&builder)

}

func main() {
	ctx := context.Background()

	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`

	conn, err := grpc.NewClient(fmt.Sprintf("%s://", builder.Scheme()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	for i := range 12 {
		log.Printf("making request: %d \n", i)

		res, err := client.GetServerAddress(ctx, &proto.GetServerAddressRequest{})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("response received: %s \n", res.GetAddress())
		time.Sleep(time.Second)
	}
}
