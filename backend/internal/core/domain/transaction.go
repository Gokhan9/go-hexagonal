package domain

import "time"

/*
Transaction, wallet içerisindeki tüm para hareketlerini (deposit, withdraw, transfer, fee vb) sistemde kayıt altına alan domain entity’dir. */
type TransactionType string

// sabitler(enums)
const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAW"
)

type Transaction struct {
	ID        string
	WalletID  string
	Amount    int64 // ! "int64" yapıldı. finansal sistemlerde yuvarlama hatalarından kaçınmak. (kuruş/cent bazında) kullanmalıyız..
	Type      TransactionType
	CreatedAt time.Time
}
