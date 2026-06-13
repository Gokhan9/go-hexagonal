package domain

import (
	"errors"
	"time"
)

// Uygulama genelinde aşağıdaki hata tanımlarını kullanacağız..
var (
	ErrorInsufficientFunds = errors.New("insufficient funds in wallet..") // Cüzdanda yetersiz bakiye
	ErrorInvalidAmount     = errors.New("amount must be greater than zero..")
)

type Wallet struct {
	ID        string
	Owner     string
	Balance   int64 // Bakiye
	Currency  string
	CreatedAt time.Time
}

// bakiyeye ekleme yapar
func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrorInvalidAmount
	}
	w.Balance += amount // w.Balance = w.Balance + amount (Mevcut bakiye(balance) üstüne amount kadar ekle.) (amount+balance(add) and balance(assign))
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
