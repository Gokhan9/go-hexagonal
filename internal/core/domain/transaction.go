package domain

import "time"

/*
*	Transaction yapısı cüzdan (wallet) içindeki para hareketlerini(İŞLEMLERİ) "KAYIT" altına almak için yazılır.
 */
type TransactionType string

// sabitler(enums)
const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "Withdraw"
)

type Transaction struct {
	ID        string
	WalletID  string
	Amount    int64 // ! "int64" yapıldı. finansal sistemlerde yuvarlama hatalarından kaçınmak. (kuruş/cent bazında) kullanmalıyız..
	Type      TransactionType
	CreatedAt time.Time
}
