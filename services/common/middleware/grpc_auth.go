package middleware

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"roomsbackend/libs/jwt"
	"roomsbackend/services/common/logger"
)

const grpcUserIDKey string = "user_id"

// UnaryAuthInterceptor проверяет JWT accessToken в gRPC метаданных.
func UnaryAuthInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("🔥 Panic в %s: %v", info.FullMethod, r)
		}
	}()

	// Публичные маршруты
	if isPublicMethod(info.FullMethod) {
		resp, err := handler(ctx, req)
		logResult(info.FullMethod, start, err)
		return resp, err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := status.Error(codes.Unauthenticated, "missing metadata")
		logResult(info.FullMethod, start, err)
		return nil, err
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
		err := status.Error(codes.Unauthenticated, "missing bearer token")
		logResult(info.FullMethod, start, err)
		return nil, err
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	userID, err := jwt.ParseToken(token)
	if err != nil {
		err = status.Error(codes.Unauthenticated, "invalid token")
		logResult(info.FullMethod, start, err)
		return nil, err
	}

	ctx = context.WithValue(ctx, grpcUserIDKey, userID)
	resp, err := handler(ctx, req)
	logResult(info.FullMethod, start, err)
	return resp, err
}

// GetUserIDFromGRPCContext извлекает user_id из gRPC контекста.
func GetUserIDFromGRPCContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(grpcUserIDKey).(string)
	return id, ok
}

// isPublicMethod возвращает true, если метод не требует авторизации.
func isPublicMethod(method string) bool {
	switch method {
	case
		"/auth.AuthService/Register",
		"/auth.AuthService/Login",
		"/auth.AuthService/Refresh":
		return true
	default:
		return false
	}
}

func logResult(method string, start time.Time, err error) {
	duration := time.Since(start).Milliseconds()

	if err != nil {
		logger.Errorf("❌ gRPC %s [error=%v] (%dms)", method, err, duration)
	} else {
		logger.Infof("✅ gRPC %s (%dms)", method, duration)
	}
}
