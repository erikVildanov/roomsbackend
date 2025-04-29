package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"roomsbackend/services/common/logger"
)

// New creates a new gRPC client connection with insecure credentials.
// Используй для внутренних сервисов внутри доверенной сети.
func New(endpoint string) *grpc.ClientConn {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(endpoint, opts)
	if err != nil {
		logger.Errorf("❌ Не удалось подключиться к gRPC %s: %v", endpoint, err)
	}
	return conn
}
