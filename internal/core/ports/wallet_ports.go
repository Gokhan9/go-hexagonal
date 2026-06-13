package ports

import (
	"context"
	"go-hexagonal/internal/core/domain"
)

// ?WalletService - Driver Port - Primary Port (Birincil)
// Uygulamanın sunduğu iş yeteneklerinin kontratıdır.
// Handler'lar (HTTP API) bu interface'i çağıracak.
type WalletService interface {
	CreateWallet(ctx context.Context, owner, currency string) (*domain.Wallet, error)
	GetWallet(ctx context.Context, id string) (*domain.Wallet, error)
	Deposit(ctx context.Context, walletID string, amount int64) error
	Withdraw(ctx context.Context, walletID string, amount int64) error
}

// ?WalletRepository - Driven Port - Secondary Port (İkincil)
// APP'in veriyi nasıl saklayacağının kontratı
// DB(postgres,redis vb..) bu interface'i implement eder.
type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	GetByID(ctx context.Context, id string) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error
}
