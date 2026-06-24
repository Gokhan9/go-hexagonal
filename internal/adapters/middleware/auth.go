package middleware

import (
	"context"
	"go-hexagonal/internal/core/domain"
	"go-hexagonal/internal/core/ports"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthenticatedUser struct {
	UserID   string
	Username string
}

/*
Authorization Header kontrolü
Bearer token parsing
JWT validation
Authenticated user’ı request context’e ekleme
*/
func AuthMiddleware(jwtSvc ports.JWTValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ! Authorization Header kontrolü
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "Authorization header required}`, http.StatusUnauthorized)
				return
			}

			// ! Bearer token parsing
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error": "Authorization format must be Bearer <token>"}`, http.StatusUnauthorized)
				return
			}

			// ! JWT validation
			claims, err := jwtSvc.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
				return
			}

			// ! Authenticated User'ı request CONTEXT'e gömmek.
			authUser := domain.AuthenticatedUser{
				UserID:   claims.UserID,
				Username: claims.Username,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, authUser)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Context'ten kullanıcı bilgisini güvenle çekmek için yardımcı fonksiyon
func GetUserFromContext(ctx context.Context) (AuthenticatedUser, error) {
	user, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return AuthenticatedUser{}, domain.ErrorUnauthorized
	}
	return user, nil
}
