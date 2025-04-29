package middleware

import (
	"context"
	"net/http"
	"strings"

	"roomsbackend/libs/jwt"
)

type contextKey string

const userIDKey contextKey = "user_id"

var publicPaths = map[string]bool{
	"/v1/auth/register": true,
	"/v1/auth/login":    true,
	"/v1/auth/refresh":  true,
}

// AuthMiddleware пропускает публичные маршруты и проверяет JWT для остальных.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwt.ParseToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID извлекает user_id из контекста.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
