package server

import (
	"net"
	"os"

	"google.golang.org/grpc"

	"roomsbackend/services/common/logger"
	"roomsbackend/services/common/middleware"
)

// GRPCRegisterFunc — функция для регистрации gRPC сервиса.
type GRPCRegisterFunc[T any] func(grpc.ServiceRegistrar, T)

// StartGRPC запускает gRPC сервер и регистрирует переданный сервис.
// T — тип хендлера (например, *AuthHandler)
func StartGRPC[T any](port string, register func(grpc.ServiceRegistrar, T), handler T) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Errorf("❌ Не удалось слушать порт %s: %v", port, err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryLoggingInterceptor, // логирование
			middleware.UnaryAuthInterceptor,    // авторизация
		),
	)
	register(grpcServer, handler)

	logger.Infof("✅ gRPC сервер запущен на %s\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Errorf("❌ gRPC сервер упал: %v", err)
	}
}
