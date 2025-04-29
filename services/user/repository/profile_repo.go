package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	userpb "roomsbackend/proto/generated/user"
)

// ProfileRepository предоставляет методы работы с профилями пользователей.
type ProfileRepository struct {
	db *sql.DB
}

// NewProfileRepository создаёт новый экземпляр ProfileRepository.
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// CreateProfile сохраняет новый профиль пользователя в БД.
func (r *ProfileRepository) CreateProfile(
	ctx context.Context,
	userID, nickname, firstName, lastName, avatarURL, bio, position, phone string,
) error {
	query := `
	INSERT INTO user_profiles (
		user_id, nickname, first_name, last_name, avatar_url, bio, position, phone_number, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		userID, nickname, firstName, lastName, avatarURL, bio, position, phone, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("ошибка создания профиля: %w", err)
	}

	return nil
}

// GetUserProfile возвращает профиль пользователя по его user_id.
func (r *ProfileRepository) GetUserProfile(ctx context.Context, userID string) (*userpb.GetUserProfileResponse, error) {
	query := `
	SELECT user_id, nickname, first_name, last_name, avatar_url, bio, position, phone_number
	FROM user_profiles
	WHERE user_id = $1
	`

	var profile userpb.GetUserProfileResponse
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserId,
		&profile.Nickname,
		&profile.FirstName,
		&profile.LastName,
		&profile.AvatarUrl,
		&profile.Bio,
		&profile.Position,
		&profile.PhoneNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("профиль не найден: %w", err)
	}

	return &profile, nil
}

// UpdateUserProfile обновляет профиль пользователя по user_id.
func (r *ProfileRepository) UpdateUserProfile(ctx context.Context, req *userpb.UpdateUserProfileRequest) error {
	query := `
	UPDATE user_profiles
	SET nickname = $1, first_name = $2, last_name = $3, avatar_url = $4,
		bio = $5, position = $6, phone_number = $7
	WHERE user_id = $8
	`

	_, err := r.db.ExecContext(ctx, query,
		req.Nickname,
		req.FirstName,
		req.LastName,
		req.AvatarUrl,
		req.Bio,
		req.Position,
		req.PhoneNumber,
		req.UserId,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления профиля: %w", err)
	}

	return nil
}

func (r *ProfileRepository) SearchUsers(ctx context.Context, query string) ([]*userpb.SearchUserResult, error) {
	sql := `
	SELECT user_id, nickname, first_name, last_name, avatar_url
	FROM user_profiles
	WHERE nickname ILIKE '%' || $1 || '%'
	   OR phone_number ILIKE '%' || $1 || '%'
	LIMIT 20
	`

	rows, err := r.db.QueryContext(ctx, sql, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска: %w", err)
	}
	defer rows.Close()

	var results []*userpb.SearchUserResult
	for rows.Next() {
		var u userpb.SearchUserResult
		if err := rows.Scan(&u.UserId, &u.Nickname, &u.FirstName, &u.LastName, &u.AvatarUrl); err != nil {
			continue
		}
		results = append(results, &u)
	}

	return results, nil
}
