package ports

import (
	"context"
	"go-hexagonal/internal/core/domain"
)

type AuditRepository interface {
	Save(ctx context.Context, log *domain.AuditLog) error
}
