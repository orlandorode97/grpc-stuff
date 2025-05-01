package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/orlandorode97/grpc-stuff/module6/internal/hello"
	"github.com/orlandorode97/grpc-stuff/module6/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
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
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "50051"
	}

	grpcServer := grpc.NewServer()

	helloService := hello.NewService(port)

	proto.RegisterHelloServiceServer(grpcServer, helloService)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
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
