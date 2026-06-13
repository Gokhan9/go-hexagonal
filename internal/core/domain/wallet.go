package domain

import (
	"errors"
	"time"
)

// Uygulama genelinde aşağıdaki hata tanımlarını kullanacağız..
var (
	ErrorInsufficientFunds = errors.New("insufficient funds in wallet..")
	ErrorInvalidAmount     = errors.New("amount must be greater than zero..")
)

type Wallet struct {
	ID        string
	Owner     string
	Balance   int64
	Currency  string
	CreatedAt time.Time
}

// bakiyeye ekleme yapar
func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrorInvalidAmount
	}
	w.Balance += amount
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
