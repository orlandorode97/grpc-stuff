package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	_ "image/png" // register the PNG format with the image package
	"io"
	"log"
	"net/http"

	"os"

	"github.com/orlandorode97/grpc-stuff/module4/exercise/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	certPool := x509.NewCertPool()
	cert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	if !certPool.AppendCertsFromPEM(cert) {
		log.Fatal("unable to append cert")
	}

	clientCert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	tlsCrendentials := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(tlsCrendentials))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := proto.NewFileUploadServiceClient(conn)

	http.HandleFunc("/download", downloadHandler(client))
	http.HandleFunc("/list", listFilesHandler(client))
	http.HandleFunc("/upload", uploadHandler(client))

	log.Printf("starting http server on: %s", ":8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal()
	}
}

func downloadHandler(client proto.FileUploadServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		imageData := make([]byte, 0)
		stream, err := client.DownloadFile(r.Context(), &proto.DownloadFileRequest{
			Name: "c526d8f0-b7d1-431c-b307-08f0d75bba29.png",
		})

		if err != nil {
			status, ok := status.FromError(err)
			if !ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Fatal(status.Message())
			http.Error(w, status.Message(), http.StatusInternalServerError)
			return
		}

		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}

				status, ok := status.FromError(err)
				if !ok {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					break
				}

				http.Error(w, status.Message(), http.StatusInternalServerError)
				break

			}

			fmt.Println("Reading content: ", len(res.Content))
			imageData = append(imageData, res.Content...)
		}

		w.Write(imageData)
	}
}

func listFilesHandler(client proto.FileUploadServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res, err := client.GetListOfFiles(r.Context(), &proto.GetListOfFilesRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(res.Files)
	}
}

func uploadHandler(client proto.FileUploadServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("/Users/orlandoromo/go/src/grpc-stuff/module4/exercise/cmd/file-upload-client/gopher.png")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		stream, err := client.UploadFile(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buffer := make([]byte, 1014)

		for {
			n, err := file.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if n == 0 {
				break
			}
			fmt.Println("Sending chunk: ", len(buffer[:n]))
			if err = stream.Send(&proto.UploadFileRequest{
				Content: buffer[:n],
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		res, err := stream.CloseAndRecv()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(res)
	}
}
