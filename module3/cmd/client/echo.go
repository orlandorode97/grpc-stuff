package main

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/orlandorode97/grpc-stuff/module3/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewStreamingServiceClient(conn)

	stream, err := client.Echo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	g, ctx := errgroup.WithContext(ctx)

	// seperate go routine to listen to the server responses
	g.Go(func() error {
		// loop for each messages from server
		for {
			res, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			// log message
			log.Printf("message received from server: %s\n", res.GetMessage())
		}

		return nil
	})

	// send messages from client
	for _, i := range []string{"#1", "#2", "#3", "#4", "#5"} {
		req := &proto.EchoRequest{Message: i}
		if err := stream.Send(req); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 2)
	}

	// close the client stream
	if err := stream.CloseSend(); err != nil {
		log.Fatal(err)
	}

	// wait for the server go routine to finish
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Println("bi-directional stream closed")

}
