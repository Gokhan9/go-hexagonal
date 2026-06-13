package dto

import (
	"go-hexagonal/internal/core/domain"
	"time"
)

// API'den gelen verileri karşılıyoruz.
type CreateWalletRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

// API'ye gönderilen verileri şekillendirme.
type WalletResponse struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json":created_at"`
}

// ?DOMAIN MODELİNİ DTO'YA ÇEVİRMEK
func ToDomainResponse(w *domain.Wallet) WalletResponse {
	return WalletResponse{
		ID:        w.ID,
		Owner:     w.Owner,
		Balance:   w.Balance,
		Currency:  w.Currency,
		CreatedAt: w.CreatedAt,
	}
}
