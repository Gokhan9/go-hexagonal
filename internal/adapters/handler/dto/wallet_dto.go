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
	Balance   float64   `json:"balance"` // UI için 10.50 olarak gösteriyoruz
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json":created_at"`
}

// ?DOMAIN MODELİNİ DTO'YA ÇEVİRMEK
func ToDomainResponse(w *domain.Wallet) WalletResponse {
	return WalletResponse{
		ID:        w.ID,
		Owner:     w.Owner,
		Balance:   float64(w.Balance) / 100, // Kuruşu TL'ye çevirip dönüyoruz
		Currency:  w.Currency,
		CreatedAt: w.CreatedAt,
	}
}
