package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/orlandorode97/grpc-stuff/module4/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {

	// public CA or public Certificate Authority
	// tslCredentials := credentials.NewTLS(&tls.Config{})
	// x509 is standard that defines the format of digital certificates used on TLS/SSL, HTTPs, VPNs etc
	// x509 starts with
	/*
				-----BEGIN CERTIFICATE-----
		    MIIC+TCCAeGgAwIBAgIJAK5N...
		    -----END CERTIFICATE-----
	*/
	certPool := x509.NewCertPool()
	cert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatal("unable to append cert")
	}

	tlsCredentials := credentials.NewClientTLSFromCert(certPool, "")

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(tlsCredentials))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	res, err := client.SayHello(context.Background(), &proto.SayHelloRequest{
		Message: "Orlando Romo",
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
