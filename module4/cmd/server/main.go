package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/orlandorode97/grpc-stuff/module4/internal/hello"
	"github.com/orlandorode97/grpc-stuff/module4/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("error running grpc server", slog.String("err", err.Error()))
		os.Exit(1)
	}

	slog.Info("shutdown grpc server gracefully")

}

func run(ctx context.Context) error {
	//loading TLS certs
	tlsCredentials, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	helloService := hello.Service{}

	proto.RegisterHelloServiceServer(grpcServer, &helloService)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			return fmt.Errorf("failted to listen on address: %w", err)
		}

		slog.Info("Starting grpc server on address", slog.String("address", lis.Addr().String()))

		if err := grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("unablet to start grpc server: %w", err)
		}
		return nil

	})

	g.Go(func() error {
		<-ctx.Done()
		grpcServer.GracefulStop()
		slog.Warn("grpc server shutdown")
		return nil
	})

	return g.Wait()
}
