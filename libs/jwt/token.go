package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"roomsbackend/services/common/logger"
)

var jwtSecret []byte

func Init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		logger.Errorf("❌ JWT_SECRET не установлен в переменных окружения")
	}
}

func GenerateAccessToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() (string, error) {
	return uuid.NewString(), nil
}

var ErrInvalidToken = errors.New("некорректный или просроченный токен")

// ParseToken проверяет токен и возвращает userID.
func ParseToken(tokenStr string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET не установлен")
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", ErrInvalidToken
	}

	// Проверяем exp вручную (безопасно)
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return "", ErrInvalidToken
		}
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	return userID, nil
}
