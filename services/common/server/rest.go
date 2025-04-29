package server

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"roomsbackend/services/common/logger"
	"roomsbackend/services/common/middleware"
)

// RESTRegisterFunc — функция для регистрации REST хендлеров через grpc-gateway.
type RESTRegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

// StartREST запускает REST API сервер через grpc-gateway.
func StartREST(port, grpcEndpoint string, register RESTRegisterFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := register(ctx, mux, grpcEndpoint, opts); err != nil {
		logger.Errorf("❌ Не удалось зарегистрировать REST хендлеры: %v", err)
	}

	// logger.Infof("✅ REST API доступен на %s\n", port)

	// ✅ Оборачиваем mux авторизацией
	logger.Infof("✅ REST API защищён JWT, доступен на %s\n", port)
	wrappedMux := middleware.AuthMiddleware(mux)

	if err := http.ListenAndServe(port, wrappedMux); err != nil {
		logger.Errorf("❌ REST сервер упал: %v", err)
	}
}
