package domain

import (
	"time"
)

var (
	ErrorInsufficientFunds = error.New("insufficient funds in wallet..")
	ErrorInvalidAmount     = error.New("amount must be greater than zero..")
)

type Wallet struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

// new wallet örneği(factory func)
func NewWallet(id, owner, currency string) *Wallet {
	return &Wallet{
		ID:        id,
		Owner:     owner,
		Balance:   0,
		Currency:  currency,
		CreatedAt: time.Now(),
	}
}

// bakiyenin çekim için yeterli mi değil mi kontrolü
func (w *Wallet) CanWithdraw(amount float64) bool {
	return w.Balance >= amount
}
