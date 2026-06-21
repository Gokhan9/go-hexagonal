package middleware

import (
	"context"
	services "go-hexagonal/internal/core/service"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthenticatedUser struct {
	UserID   string
	Username string
}

func AuthMiddleware(jwtSvc *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "Authorization header required}`, http.StatusUnauthorized)
				return
			}

			//
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error": "Authorization format must be Bearer <token>"}`, http.StatusUnauthorized)
				return
			}

			//
			claims, err := jwtSvc.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
				return
			}

			// User'ı CONTEXT içine gömmek.
			authUser := AuthenticatedUser{
				UserID:   claims.UserID,
				Username: claims.Username,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, authUser)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
