package ports

import "go-hexagonal/internal/core/domain"

type JWTValidator interface {
	ValidateToken(tokenStr string) (*domain.UserClaims, error)
}
