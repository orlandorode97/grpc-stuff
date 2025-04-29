package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/orlandorode97/grpc-stuff/module5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// WithChainUnaryInterceptor will chain multiple uninary interceptor, the first interceptor will be the outer most, while the last interceptor will be
		// inner most of the call.
		grpc.WithChainUnaryInterceptor(
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				start := time.Now()

				err := invoker(ctx, method, req, reply, cc, opts...)

				log.Printf("request %s took %s\n", method, time.Since(start))
				return err
			},
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				log.Printf("sending request: %s\n\n", method)

				err := invoker(ctx, method, req, reply, cc, opts...)

				log.Printf("This was the req: %v\n\n", req)
				log.Printf("This is the reply: %v\n\n", reply)
				log.Printf("This is the client connection: %v\n\n", cc.GetState().String())

				log.Printf("response received from server: %s\n", method)

				return err
			}),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	res, err := client.SayHello(context.Background(), &proto.SayHelloRequest{
		Name: "Orlandp",
	})

	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			log.Fatal("unable to parse error", status.Err().Error())
		}

		log.Fatal(status.Message())
	}

	fmt.Printf("response received: %s\n", res.Message)
}
