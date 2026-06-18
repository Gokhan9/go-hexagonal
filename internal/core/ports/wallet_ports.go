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
	GetWallet(ctx context.Context, id string) (*domain.Wallet, error)

	// ! Adım 1: "ports.WalletService" interface'indeki "Deposit" ve "Withdraw" metod imzalarını güncelleyeceğiz ("idempotencyKey" parametresini ekleyeceğiz).
	// ! "idempotencyKey" parametreleri  YENİ EKLENDİ.
	Deposit(ctx context.Context, idempotencyKey string, walletID string, amount int64) error
	Withdraw(ctx context.Context, idempotencyKey string, walletID string, amount int64) error
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

	// ! YENİ EKLENDİ.
	SaveTransaction(ctx context.Context, tn *domain.Transaction) error                             // İşlem kaydını kalıcı hale getirir, save'ler
	GetTransactionsByWalletID(ctx context.Context, walletID string) ([]*domain.Transaction, error) // ID'ye göre cüzdana ait tüm geçmiş hareketleri getirir.

}

type UserService interface {
	Register(ctx context.Context, username, password string) (*domain.User, error) // KAYIT OL - New User Creating(Username+Password)-Password Hash(SetPassword)-User DB Kaydedilir.
	Login(ctx context.Context, username, password string) (string, error)          // GİRİŞ YAP - User'ı bulur(token verir.)-Şifreyi kontrol(CheckPassword)-Success(JWT token/session token(string))-Fail(Login Err)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
