package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// RefreshTokenRepository управляет таблицей refresh_tokens.
type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository создает новый экземпляр RefreshTokenRepository.
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Save сохраняет refresh token в базу данных.
func (r *RefreshTokenRepository) Save(ctx context.Context, token, userID string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, token, userID, expiresAt)
	return err
}

// Delete удаляет конкретный refresh token из базы.
func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM refresh_tokens WHERE token = $1
	`, token)
	return err
}

// Exists проверяет, существует ли токен.
func (r *RefreshTokenRepository) Exists(ctx context.Context, token string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM refresh_tokens WHERE token = $1)
	`, token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке токена: %w", err)
	}
	return exists, nil
}
