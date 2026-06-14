package domain

import "errors"

// Uygulama genelinde aşağıdaki hata tanımlarını kullanacağız..
var (
	ErrorInsufficientFunds = errors.New("insufficient funds in wallet..")     // Cüzdanda Yetersiz Bakiye bakiye
	ErrorInvalidAmount     = errors.New("amount must be greater than zero..") // Geçersiz Miktar
)
