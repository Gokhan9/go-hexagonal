package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey []byte
	issuer    string
	expiry    time.Duration
}

type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, expiry time.Duration) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		issuer:    "wallet-api",
		expiry:    expiry,
	}
}
