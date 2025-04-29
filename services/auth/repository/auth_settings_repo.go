package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SessionRepository управляет настройками сессий пользователя (например TTL).
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository создаёт SessionRepository.
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// UpdateSessionTTL обновляет session_ttl в users.
func (r *SessionRepository) UpdateSessionTTL(ctx context.Context, userID, ttl string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET session_ttl = $1 WHERE id = $2
	`, ttl, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления TTL: %w", err)
	}
	return nil
}

// GetSessionTTL получает session_ttl и преобразует в time.Duration.
func (r *SessionRepository) GetSessionTTL(ctx context.Context, userID string) (time.Duration, error) {
	var raw string
	err := r.db.QueryRowContext(ctx, `
		SELECT session_ttl FROM users WHERE id = $1
	`, userID).Scan(&raw)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения TTL из БД: %w", err)
	}

	ttl, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("невалидный session_ttl: %w", err)
	}

	return ttl, nil
}
