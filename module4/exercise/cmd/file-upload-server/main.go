package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"

	"github.com/orlandorode97/grpc-stuff/module4/exercise/internal/uploading"
	"github.com/orlandorode97/grpc-stuff/module4/exercise/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	serverCert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal(err)
		return
	}

	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatal(err)
		return
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatal("unable to append ca cert")
		return
	}

	tlsCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	service := uploading.New()

	proto.RegisterFileUploadServiceServer(grpcServer, service)

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting grpc server on: %v\n", listen.Addr().String())

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
