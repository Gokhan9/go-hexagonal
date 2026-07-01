package ports

import (
	"context"
	"go-hexagonal/internal/core/domain"
)

type UserService interface {
	Register(ctx context.Context, username, password string) (*domain.User, error) // KAYIT OL - New User Creating(Username+Password)-Password Hash(SetPassword)-User DB Kaydedilir.
	Login(ctx context.Context, username, password string) (string, error)          // GİRİŞ YAP - User'ı bulur(token verir.)-Şifreyi kontrol(CheckPassword)-Success(JWT token/session token(string))-Fail(Login Err)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
