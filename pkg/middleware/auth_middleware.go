package middleware

import (
	"context"
	"net/http"
	"strings"
	"x-clone/internal/service"
)

type ContextKey string

const (
	UserIDKey ContextKey = "userID"
)

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "you are not authorised", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authService.ValidateAccessToken(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userID := int(claims["user_id"].(float64))
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
