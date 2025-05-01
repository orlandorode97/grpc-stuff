package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/orlandorode97/grpc-stuff/module5/internal/auth"
	"github.com/orlandorode97/grpc-stuff/module5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type tokenIssuer interface {
	IssueToken(ctx context.Context, userID string) (string, error)
}
type jwtCredentials struct {
	tokenIssuer tokenIssuer
}

func (j *jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	info, ok := credentials.RequestInfoFromContext(ctx)
	if !ok || info.Method != proto.HelloService_Protected_FullMethodName {
		return nil, nil
	}

	token, err := j.tokenIssuer.IssueToken(ctx, "user-id-12345")
	if err != nil {
		return nil, err
	}

	return map[string]string{"authorization": token}, nil
}
func (j *jwtCredentials) RequireTransportSecurity() bool {
	return false
}

func main() {
	ctx := context.Background()
	jwt, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET is reqired")
	}
	auth, err := auth.NewService(jwt)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&jwtCredentials{tokenIssuer: auth}),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewHelloServiceClient(conn)

	res, err := client.Protected(ctx, &proto.ProtectedRequest{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("response received: %s\n", res.GetUserId())
}
