package repository

import (
	"context"
	"errors"
	"go-hexagonal/internal/core/domain"
	"sync"
)

/*
  Kodun Açıklaması:
   * Thread-Safety: sync.RWMutex kullanarak haritanın (map) eşzamanlı (concurrent) okuma ve yazma işlemlerinde bozulmasını engelleriz.
   * Port Uyumluluğu: internal/core/ports/wallet_ports.go içinde tanımladığımız WalletRepository interface'indeki tüm metodları (Create, GetByID, Update) somutlaştırırız.
   * Veri Saklama: Veriler uygulama çalıştığı sürece bir map içinde tutulur; uygulama kapandığında veriler silinir. Bu, geliştirme aşamasında hızlı prototipleme sağlar.

*/

// "MemoryWalletRepository", WalletRepository interface'ini memory üzerinden implement eder.
type MemoryWalletRepository struct {
	wallets map[string]*domain.Wallet
	mu      sync.RWMutex
}

// "NewMemoryWalletRepository", yeni bir memory deposu create eder.
func NewMemoryWalletRepository() *MemoryWalletRepository {
	return &MemoryWalletRepository{
		wallets: make(map[string]*domain.Wallet),
	}
}

// "Create" ile yeni bir cüzdanı memory'e kaydeder.
func (r *MemoryWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	r.mu.Lock()         // → Aynı anda farklı goroutine'lerin "wallets" map'ine erişimini kısıtlar.
	defer r.mu.Unlock() // → function bitince kilidi açar. (RACE CONDITION önlemek.)

	// → "wallet.ID" ile "map" içinde arama yapıyor
	if _, exists := r.wallets[wallet.ID]; exists {
		return errors.New("wallet already exists.")
	}

	r.wallets[wallet.ID] = wallet // ? Yeni "wallet", "Map" içine eklenir.
	return nil
}

// GetByID'ye göre cüzdanı getir
func (r *MemoryWalletRepository) GetByID(ctx context.Context, id string) (*domain.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	wallet, exists := r.wallets[id]
	if !exists {
		return nil, errors.New("wallet not found..")
	}

	return wallet, nil
}

// Update ile mevcut olan wallet'i günceller..
func (r *MemoryWalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.wallets[wallet.ID]; !exists {
		return errors.New("wallet not found")
	}

	r.wallets[wallet.ID] = wallet
	return nil
}
