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
	CreatedAt time.Time `json:"created_at"`
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

/*
TODO: UI/API dış dünyada parayı "float64" (örneğin 10.50 TL) olarak gönderiyor. Ancak domain katmanında kuruş hassasiyeti ve hassas finansal hesaplamalar için bunu "int64" (kuruş/sent bazlı: 1050 kuruş) olarak tutuyoruz.
*Para yatırma ve çekme requestlerini karşılamak için tek bir ortak request modeli.
*/
type TransactionRequest struct {
	Amount float64 `json:"amount"`
}

/*
TODO: Gelen "float" tutarı, kuruş(int64) cinsine çeviren yardımcı bir metod.
*Finansal işlemlerde "float" yuvarlama hatalarını önlemek için "0.5" ekleyerek "cast" ediyoruz.
*/
func (r *TransactionRequest) ToCents() int64 {
	return int64((r.Amount * 100) + 0.5)
}
