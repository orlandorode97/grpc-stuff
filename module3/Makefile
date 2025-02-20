get-protoc-plugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/streaming.proto
run-server:
	echo "runnning-server"
	go run cmd/streaming/main.go
run-client:
	echo "running client"
	go run cmd/client/main.go

