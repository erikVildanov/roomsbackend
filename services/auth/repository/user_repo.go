package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, login, password string) (string, error) {
	userID := uuid.New().String()

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	query := `INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3)`
	_, err = r.db.ExecContext(ctx, query, userID, login, string(hashed))
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения пользователя: %w", err)
	}

	return userID, nil
}

func (r *UserRepository) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
	var userID string
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return "", fmt.Errorf("refresh token не найден или ошибка: %w", err)
	}

	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("refresh token истёк")
	}

	return userID, nil
}

// CheckUserCredentials проверяет логин и пароль пользователя.
// Возвращает userID при успешной проверке или ошибку при несовпадении данных.
func (r *UserRepository) CheckUserCredentials(ctx context.Context, login, password string) (string, error) {
	var userID string
	var hashedPassword string

	query := `SELECT id, password_hash FROM users WHERE login = $1`
	err := r.db.QueryRowContext(ctx, query, login).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", fmt.Errorf("пользователь не найден: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("неверный пароль")
	}

	return userID, nil
}
