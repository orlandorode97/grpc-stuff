get-protoc-plugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/hello.proto
run-server:
	PORT=50051 go run cmd/server/main.go &
	PORT=50052 go run cmd/server/main.go &
	PORT=50053 go run cmd/server/main.go
