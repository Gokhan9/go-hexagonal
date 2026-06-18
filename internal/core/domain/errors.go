package domain

import "errors"

// Uygulama genelinde aşağıdaki hata tanımlarını kullanacağız..
var (
	ErrorInsufficientFunds    = errors.New("Insufficient funds in wallet.")      // Cüzdanda Yetersiz Bakiye
	ErrorInvalidAmount        = errors.New("Amount must be greater than 0.....") // Geçersiz Miktar
	ErrConcurrentModification = errors.New("Eşzamanlı değişiklik hatası..")

	// Auth Errors
	ErrorUserAlreadyExists  = errors.New("Username is already taken")    // ilgili username'in önceden alınmış olması
	ErrorInvalidCredentials = errors.New("Invalid username or password") // login sırasında username/e-posta eşleşmemesi
	ErrorUnauthorized       = errors.New("Unauthorized access")          // kimlik doğrulama(authentication) başarısız(401).
)
