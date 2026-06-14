package services

import (
	"context"
	"go-hexagonal/internal/core/domain"
	"go-hexagonal/internal/core/ports"

	"github.com/google/uuid"
)

// dependency inject
type walletService struct {
	repo ports.WalletRepository
}

func NewWalletService(repo ports.WalletRepository) ports.WalletService {
	return &walletService{
		repo: repo,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, owner, currency string) (*domain.Wallet, error) {
	// factory methodu ile yeni bir wallet oluşturmak..
	wallet := &domain.Wallet{
		ID:       uuid.NewString(),
		Owner:    owner,
		Balance:  0,
		Currency: currency,
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) GetWallet(ctx context.Context, id string) (*domain.Wallet, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *walletService) Deposit(ctx context.Context, walletID string, amount int64) error {

	if amount <= 0 {
		return domain.ErrorInvalidAmount // <-- Guard Clause(Koruyucu Koşul)
	}

	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	// İş modelini domain modeli üzerindeki metodla işletmek.
	if err := wallet.Deposit(amount); err != nil {
		return err
	}

	return s.repo.Update(ctx, wallet)
}

func (s *walletService) Withdraw(ctx context.Context, walletID string, amount int64) error {

	if amount <= 0 {
		return domain.ErrorInsufficientFunds // <-- Guard Clause
	}

	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if err := wallet.Withdraw(amount); err != nil {
		return err
	}

	return s.repo.Update(ctx, wallet)
}
