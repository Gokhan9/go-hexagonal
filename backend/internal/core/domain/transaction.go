package domain

import "time"

type TransactionType string   // Transaction, wallet içerisindeki tüm para hareketlerini (deposit, withdraw, transfer, fee vb) sistemde kayıt altına alan domain entity’dir.
type TransactionStatus string // Transaction İşlemi başladığında PENDING, başarı durumunda COMPLETED, hata durumunda FAILED statüsüne geçirmek.

// SABİTLER(enums)
const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAW"
)

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusCompleted TransactionStatus = "COMPLETED"
	StatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID        string
	WalletID  string
	Amount    int64 // ! "int64" yapıldı. finansal sistemlerde yuvarlama hatalarından kaçınmak. (kuruş/cent bazında) kullanmalıyız..
	Type      TransactionType
	Status    TransactionStatus
	CreatedAt time.Time
}
