package domain

import "errors"

// Uygulama genelinde aşağıdaki hata tanımlarını kullanacağız..
var (
	ErrorInsufficientFunds    = errors.New("Insufficient funds in wallet.")  // Cüzdanda Yetersiz Bakiye
	ErrorInvalidAmount        = errors.New("Amount must be greater than 0.") // Geçersiz Miktar
	ErrConcurrentModification = errors.New("Eşzamanlı değişiklik hatası.")

	// Auth Errors
	ErrorUserAlreadyExists  = errors.New("Username is already taken")    // İlgili username'in önceden alınmış olması
	ErrorInvalidCredentials = errors.New("Invalid Username or Password") // Login sırasında username/e-posta eşleşmemesi
	ErrorUnauthorized       = errors.New("Unauthorized Access")          // Kimlik doğrulama(authentication) başarısız(401).

	// wallet error
	ErrorWalletNotFound = errors.New("Wallet not found..")

	ErrorSelfTransfer        = errors.New("Self transfer is not allowed.")  // Kendi cüzdanına transfer yapılamaz.
	ErrorWalletAlreadyClosed = errors.New("Wallet is already closed.")      // Wallet kapatma
	ErrorWalletNotEmptied    = errors.New("Wallet must be empty to close.") // Wallet kapatmak için boş olmalı.
)
