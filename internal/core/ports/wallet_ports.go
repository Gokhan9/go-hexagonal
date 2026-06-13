package ports

import (
	"context"
	"go-hexagonal/internal/core/domain"
)

// WalletRepository-Driven Port(İkincil)
// APP'in veriyi nasıl saklayacağının kontratı
// DB(postgres,redis vb..) bu interface'i implement eder.
type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	GetByID(ctx context.Context, id string) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error
}

// WalletService - Driving Port (Birincil Liman)
// Uygulamanın sunduğu iş yeteneklerinin kontratıdır.
// Handler'lar (HTTP API) bu interface'i çağıracak.
