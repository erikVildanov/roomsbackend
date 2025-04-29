package handler

import (
	"context"

	"roomsbackend/proto/generated/user"
	"roomsbackend/services/common/logger"
	"roomsbackend/services/user/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler реализует gRPC-интерфейс UserService.
type UserHandler struct {
	user.UnimplementedUserServiceServer
	Repo *repository.ProfileRepository
}

// CreateUserProfile обрабатывает создание пользовательского профиля.
// Вызывается после регистрации из auth-сервиса.
// Сохраняет профиль в БД через репозиторий.
func (h *UserHandler) CreateUserProfile(ctx context.Context, req *user.CreateUserProfileRequest) (*user.CreateUserProfileResponse, error) {
	err := h.Repo.CreateProfile(ctx,
		req.UserId,
		req.Nickname,
		req.FirstName,
		req.LastName,
		req.AvatarUrl,
		req.Bio,
		req.Position,
		req.PhoneNumber,
	)

	if err != nil {
		logger.Errorf("❌ Ошибка создания профиля:", err)
		return nil, err
	}

	return &user.CreateUserProfileResponse{
		Status: "ok",
	}, nil
}

// GetUserProfile возвращает профиль пользователя по ID.
func (h *UserHandler) GetUserProfile(ctx context.Context, req *user.GetUserProfileRequest) (*user.GetUserProfileResponse, error) {
	profile, err := h.Repo.GetUserProfile(ctx, req.UserId)
	if err != nil {
		logger.Errorf("❌ Ошибка получения профиля: %v", err)
		return nil, err
	}
	return profile, nil
}

// UpdateUserProfile обновляет профиль пользователя.
func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *user.UpdateUserProfileRequest) (*user.UpdateUserProfileResponse, error) {
	err := h.Repo.UpdateUserProfile(ctx, req)
	if err != nil {
		logger.Errorf("❌ Ошибка обновления профиля: %v", err)
		return nil, err
	}
	return &user.UpdateUserProfileResponse{
		Status: "ok",
	}, nil
}

func (h *UserHandler) SearchUsers(ctx context.Context, req *user.SearchUsersRequest) (*user.SearchUsersResponse, error) {
	results, err := h.Repo.SearchUsers(ctx, req.Query)
	if err != nil {
		logger.Errorf("❌ Ошибка поиска: %v", err)
		return nil, status.Error(codes.Internal, "ошибка поиска")
	}
	return &user.SearchUsersResponse{Users: results}, nil
}
