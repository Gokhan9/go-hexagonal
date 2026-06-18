package services

import (
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

// JWT içine koyulacak payload(veri) kısmı.
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
