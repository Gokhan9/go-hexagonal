package ports

import (
	"context"
	"go-hexagonal/internal/core/domain"
)

// ? WalletService - Driver Port - Primary Port (Birincil)
// Uygulamanın sunduğu iş yeteneklerinin kontratıdır.
// Handler'lar (HTTP API) bu interface'i çağıracak.
type WalletService interface {
	CreateWallet(ctx context.Context, owner, currency string) (*domain.Wallet, error)
	GetWallet(ctx context.Context, userID, id string) (*domain.Wallet, error)

	// ! Adım 1: "ports.WalletService" interface'indeki "Deposit" ve "Withdraw" metod imzalarını güncelleyeceğiz ("idempotencyKey" parametresini ekleyeceğiz).
	// ! "idempotencyKey" parametreleri  YENİ EKLENDİ.
	Deposit(ctx context.Context, idempotencyKey string, walletID string, userID string, TransactionID string, amount int64) error
	Withdraw(ctx context.Context, idempotencyKey string, walletID string, userID string, TransactionID string, amount int64) error
	GetTransactions(ctx context.Context, walletID string) ([]*domain.Transaction, error)
}

// ? WalletRepository - Driven Port - Secondary Port (İkincil)
// APP'in veriyi nasıl saklayacağının kontratı
// DB(postgres,redis vb..) bu interface'i implement eder.
type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	GetByID(ctx context.Context, id string) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error

	GetIdempotencyRecord(ctx context.Context, key string) (*domain.IdempotencyRecord, error) // "KEY"'in önceden kullanılıp, kullanılmadığını kontrol edeceğiz.
	SaveIdempotencyRecord(ctx context.Context, record *domain.IdempotencyRecord) error       // Yeni işlemi "KEY" ile kaydedeceğiz.

	SaveTransaction(ctx context.Context, tn *domain.Transaction) error                             // İşlem kaydını kalıcı hale getirir, save'ler
	GetTransactionsByWalletID(ctx context.Context, walletID string) ([]*domain.Transaction, error) // ID'ye göre cüzdana ait tüm geçmiş hareketleri getirir.

	BeginTx(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
