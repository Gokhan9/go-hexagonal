package domain

import "github.com/golang-jwt/jwt/v5"

// JWT Token içine koyulacak payload(veri) kısmı.
type UserClaims struct {
	UserID               string `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // Embedded. Standart JWT alanlarını ekler. (iss-issuer, exp-expiration(bitiş zamanı), iat-issued at(oluşturulma zamanı), sub(subject))
}
