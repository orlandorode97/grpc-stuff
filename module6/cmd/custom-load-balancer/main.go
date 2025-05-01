package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/orlandorode97/grpc-stuff/module6/internal/loadbalancer"
	"github.com/orlandorode97/grpc-stuff/module6/internal/resolve"
	"github.com/orlandorode97/grpc-stuff/module6/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)

var builder resolve.Builder

func init() {
	builder := resolve.Builder{}
	resolver.Register(&builder)

}

func main() {
	ctx := context.Background()

	groups := map[string]string{
		"group-a": "localhost:50051",
		"group-b": "localhost:50052",
	}

	lbBuilder := loadbalancer.NewBuilder(groups, "localhost:50053")

	balancer.Register(lbBuilder)

	serviceConfig := `{"loadBalancingPolicy": "ab_testing"}`

	conn, err := grpc.NewClient(fmt.Sprintf("%s://", builder.Scheme()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	time.Sleep(time.Second)

	client := proto.NewHelloServiceClient(conn)

	for _, group := range []string{"group-a", "group-b", "group-c"} {
		log.Printf("making request for group: %q \n", group)
		res, err := client.GetServerAddress(
			metadata.AppendToOutgoingContext(ctx, "user-group", group),
			&proto.GetServerAddressRequest{},
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("response receive for group %q  from server: %s \n", group, res.GetAddress())
	}
}
