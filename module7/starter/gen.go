package starter

// gRPC Layer
//go:generate mockgen -destination=internal/mocks/chat/mock_chat_service.gen.go -package=chat_mock github.com/orlandorode97/grpc-stuff/module7/starter/internal/transport/grpc ChatService
//go:generate mockgen -destination=internal/mocks/grpc/mock_grpc_stream.gen.go -package=grpc_mock github.com/orlandorode97/grpc-stuff/module7/starter/proto ChatService_SubscribeServer

// Service Layer
//go:generate mockgen -destination=internal/mocks/chat/mock_id_generator.gen.go -package=chat_mock github.com/orlandorode97/grpc-stuff/module7/starter/internal/chat IDGenerator
//go:generate mockgen -destination=internal/mocks/store/mock_store.gen.go -package=store_mock github.com/orlandorode97/grpc-stuff/module7/starter/internal/chat Store
