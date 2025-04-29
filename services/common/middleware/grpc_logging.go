package middleware

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"roomsbackend/services/common/logger"
)

// UnaryLoggingInterceptor логирует входные параметры, ошибки и время выполнения.
func UnaryLoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	// Лог запроса
	payload, _ := json.Marshal(req)
	logger.Debugf("➡️ gRPC %s | req: %s", info.FullMethod, payload)

	// Вызов метода
	resp, err := handler(ctx, req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		st, _ := status.FromError(err)
		logger.Errorf("❌ gRPC %s | error: %s (%dms)", info.FullMethod, st.Message(), duration)
		return nil, err
	}

	// Лог ответа
	respPayload, _ := json.Marshal(resp)
	logger.Debugf("⬅️ gRPC %s | resp: %s (%dms)", info.FullMethod, respPayload, duration)

	return resp, nil
}
