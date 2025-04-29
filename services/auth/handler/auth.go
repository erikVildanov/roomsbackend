package handler

import (
	"context"
	"time"

	"roomsbackend/proto/generated/auth"
	userpb "roomsbackend/proto/generated/user"
	"roomsbackend/services/auth/repository"
	"roomsbackend/services/common/logger"

	"roomsbackend/libs/jwt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthHandler обрабатывает gRPC-запросы AuthService.
type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	Repo        *repository.UserRepository
	RefreshRepo *repository.RefreshTokenRepository
	SessionRepo *repository.SessionRepository
	UserGRPC    userpb.UserServiceClient // gRPC-клиент user-сервиса
}

func (h *AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	userID, err := h.Repo.CreateUser(ctx, req.Login, req.Password)
	if err != nil {
		logger.Errorf("❌ ошибка регистрации: %v", err)
		return nil, err
	}

	_, err = h.UserGRPC.CreateUserProfile(ctx, &userpb.CreateUserProfileRequest{
		UserId:      userID,
		Nickname:    req.Nickname,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		AvatarUrl:   req.AvatarUrl,
		Bio:         "", // пока пусто
		Position:    "",
		PhoneNumber: "",
	})
	if err != nil {
		logger.Errorf("❌ Ошибка при попытке создать профиль пользователя: %v", err)
		return nil, err
	}

	ttl := time.Hour * 24 * 7 // 7 дней по умолчанию

	accessToken, err := jwt.GenerateAccessToken(userID, ttl)
	if err != nil {
		logger.Errorf("❌ Ошибка генерации access token: %v", err)
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		logger.Errorf("❌ Ошибка генерации refresh token: %v", err)
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	err = h.RefreshRepo.Save(ctx, refreshToken, userID, expiresAt)
	if err != nil {
		logger.Errorf("❌ Не удалось сохранить новый refresh токен: %v", err)
		return nil, err
	}

	return &auth.RegisterResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	userID, err := h.Repo.GetUserIDByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		logger.Errorf("❌ Ошибка refresh токена: %v", err)
		return nil, err
	}

	ttl := time.Hour * 24 * 7 // 7 дней по умолчанию

	if userTTL, err := h.SessionRepo.GetSessionTTL(ctx, userID); err == nil {
		ttl = userTTL
	}

	accessToken, err := jwt.GenerateAccessToken(userID, ttl)
	if err != nil {
		logger.Errorf("❌ Ошибка генерации access token: %v", err)
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		logger.Errorf("❌ Ошибка генерации refresh token: %v", err)
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	err = h.RefreshRepo.Save(ctx, refreshToken, userID, expiresAt)
	if err != nil {
		logger.Errorf("❌ Не удалось сохранить новый refresh токен: %v", err)
		return nil, err
	}

	return &auth.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login обрабатывает вход пользователя по логину и паролю.
// При успешной авторизации возвращает новый access и refresh токены.
func (h *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	userID, err := h.Repo.CheckUserCredentials(ctx, req.Login, req.Password)
	if err != nil {
		logger.Errorf("❌ Ошибка логина: %v", err)
		return nil, err
	}

	ttl := time.Hour * 24 * 7 // 7 дней по умолчанию

	if userTTL, err := h.SessionRepo.GetSessionTTL(ctx, userID); err == nil {
		ttl = userTTL
	}

	accessToken, err := jwt.GenerateAccessToken(userID, ttl)
	if err != nil {
		logger.Errorf("❌ Ошибка генерации access токена: %v", err)
		return nil, err
	}

	refreshToken, _ := jwt.GenerateRefreshToken()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	if err := h.RefreshRepo.Save(ctx, refreshToken, userID, expiresAt); err != nil {
		logger.Errorf("❌ Ошибка сохранения refresh токена при логине: %v", err)
		return nil, err
	}

	return &auth.LoginResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	err := h.RefreshRepo.Delete(ctx, req.RefreshToken)
	if err != nil {
		logger.Errorf("❌ Ошибка при logout: %v", err)
		return nil, status.Error(codes.Internal, "не удалось завершить сессию")
	}
	return &auth.LogoutResponse{Status: "ok"}, nil
}
