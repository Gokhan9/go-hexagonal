package services

import (
	"context"
	"go-hexagonal/internal/core/domain"
	"go-hexagonal/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo   ports.UserRepository
	jwtService *JWTService
}

func NewUserService(repo ports.UserRepository, jwtSvc *JWTService) ports.UserService {
	return &UserService{
		userRepo:   repo,
		jwtService: jwtSvc,
	}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*domain.User, error) {

	// Username Check
	// && operatöründe ilk koşul "false" ise, Go ikinci koşula hiç bakmaz (çünkü sonuç zaten false olacaktır).
	exist, err := s.userRepo.GetByUsername(ctx, username)
	if err == nil && exist != nil {
		return nil, domain.ErrorUserAlreadyExists
	}

	// New User Creating
	user := &domain.User{
		ID:        uuid.NewString(),
		Username:  username,
		CreatedAt: time.Now(),
	}

	// Password Hash and Save
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (string, error) {

	// || operatöründe ilk koşul "true" ise, Go ikinci koşula hiç bakmaz (çünkü sonuç zaten true olacaktır).
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil || user == nil {
		return "", domain.ErrorInvalidCredentials
	}

	// Hash Doğrula
	if !user.CheckPassword(password) {
		return "", domain.ErrorInvalidCredentials
	}

	// Token üret
	return s.jwtService.GenerateToken(user.ID, user.Username)
}
