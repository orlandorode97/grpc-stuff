package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
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
	serverCert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		return fmt.Errorf("unable to load x509 tls certs: %w", err)
	}
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return fmt.Errorf("unable to open CA cert: %w", err)
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(caCert) {
		return errors.New("unable to append certs")
	}

	tlsCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

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
