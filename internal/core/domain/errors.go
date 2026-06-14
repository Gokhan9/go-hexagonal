package domain

import "errors"

// Uygulama genelinde aşağıdaki hata tanımlarını kullanacağız..
var (
	ErrorInsufficientFunds = errors.New("Insufficient funds in wallet.")      // Cüzdanda Yetersiz Bakiye
	ErrorInvalidAmount     = errors.New("Amount must be greater than 0.....") // Geçersiz Miktar
)
