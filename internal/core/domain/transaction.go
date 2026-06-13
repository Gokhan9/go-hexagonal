package domain

import "time"

/*
*	Transaction yapısı cüzdan (wallet) içindeki para hareketlerini KAYIT altına almak için yazılır.
 */
type TransactionType string

// sabitler(enums)
const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
)

type Transaction struct {
	ID        string
	WalletID  string
	Amount    float64
	Type      TransactionType
	CreatedAt time.Time
}
