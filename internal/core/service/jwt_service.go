package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
1. Token üretme (muhtemelen başka fonksiyonda)
-Kullanıcı login olur
-UserClaims oluşturulur
-JWT imzalanır

2. Token doğrulama
-Secret key ile imza kontrol edilir
-Expiration kontrol edilir
*/

// JWT Üretmek ve Doğrulamak.
type JWTService struct {
	secretKey []byte        // Tokenları imzalamak ve doğrulamak için unique key. JWT kütüphanesi "[]byte" ister.
	issuer    string        // Token'ı üreten bilgisi(örn:wallet-api).
	expiry    time.Duration // Token süresi
}

// JWT Token içine koyulacak payload(veri) kısmı.
type UserClaims struct {
	UserID               string `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // Embedded. Standart JWT alanlarını ekler. (iss-issuer, exp-expiration(bitiş zamanı), iat-issued at(oluşturulma zamanı), sub(subject))
}

/*
- "secret" string alıyor, byte'a çevirip saklıyor.
- "issuer" sabit olarak, "wallet-api"
- "expirt" dışarıdan geliyor.
*/
func NewJWTService(secret string, expiry time.Duration) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    "wallet-api",
		expiry:    expiry,
	}
}

// Her user için JWT Token üretir.
func (s *JWTService) GenerateToken(UserID, Username string) (string, error) {

	// JWT Token içine koyulacak veriyi(UserClaims) oluşturduk.
	claims := &UserClaims{
		UserID:   UserID,   // TOKEN İÇİNDE UserID saklıyoruz.
		Username: Username, // TOKEN İÇİNDE Username saklıyoruz.
		// JWT Standartlarına ait default alanlar.
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,                                     // TOKEN ÜRETEN BİLGİSİ(wallet-api).
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // TOKEN OLUŞTURULMA ZAMANI.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)), // TOKEN GEÇERLİLİK SÜRESİ.
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims) // "ES256", algoritmasıyla JWT instance create.
	return token.SignedString(s.secretKey)                     // Token, secretKey ile imzalayıp string olarak döneriz.
}

// Token validate
func (s *JWTService) ValidateToken(tokenStr string) (*UserClaims, error) {

	/*
		token, err := jwt.ParseWithClaims(tokenStr) - TokenStr parse(çözümleme) ediyor.
		UserClaims{} - Token içindeki payload(verilerin(userid,username)) dönüştürüleceği struct.
		func(token *jwt.Token) (interface{}, error) - İmza doğrulama sırasında çalışan rollback.
	*/
	token, err := jwt.ParseWithClaims(tokenStr, UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Token'ın "HMAC" ile imzalandığını kontrol eder, farklı algoritma ile oluşturulan tokenları engeller.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method..")
		}
		return s.secretKey, nil // İmzayı doğrulamak için secretKey dönülür.
	})

	// Parse veya doğrulama sonrası hata olursa "err" döner.
	if err != nil {
		return nil, err
	}

	// Token içinde ki Claims'leri, UserClaims'e dönüştür ayrıca token geçerli mi(valid) kontrol et.
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil // User Claims bilgilerini döndür
	}

	return nil, errors.New("invalid token claims..") // token geçersizse hata dön.
}
