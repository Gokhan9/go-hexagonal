package domain

import (
	"time"
)

type Wallet struct {
	ID        string
	OwnerID   string // UPTADE: Artık UserID'ye bağlı OwnerID
	Balance   int64  // Bakiye
	Currency  string // Para Birimi
	CreatedAt time.Time

	Version int
}

// bakiyeye ekleme yapar
func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrorInvalidAmount
	}
	w.Balance += amount // w.Balance = Balance + amount (Mevcut bakiye(balance) üstüne amount kadar ekle.) (amount+balance(add) and balance(assign))
	return nil
}

// bakiyeden düşer.
func (w *Wallet) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrorInvalidAmount
	}
	if w.Balance < amount {
		return ErrorInsufficientFunds
	}
	w.Balance -= amount
	return nil
}
