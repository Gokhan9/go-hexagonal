package service

import (
	"context"
	"errors"
	"fmt"
	"go-hexagonal/internal/core/domain"
	"go-hexagonal/internal/core/ports"
	"time"

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
		OwnerID:  owner, // UPDATE: OwnerID
		Balance:  0,
		Currency: currency,
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) GetWallet(ctx context.Context, userID, id string) (*domain.Wallet, error) {

	wallet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// * 2- YETKI KONTROLÜ: Bu wallet istek atan user'a mı ait?
	if wallet.OwnerID != userID {
		return nil, domain.ErrorUnauthorized
	}

	return wallet, nil
}

// ! 2. ADIM : WalletService, "Deposit ve Withdraw" Güncellemesi
func (s *walletService) Deposit(ctx context.Context, idempotencyKey string, walletID string, userID string, transactionID string, amount int64) error {

	if amount <= 0 {
		return domain.ErrorInvalidAmount // Guard Clause
	}

	// ! 1. Idempotency Kontrolü
	if idempotencyKey != "" {

		// ! 1. Sorgulama (Check): Eğer "idempotencyKey" boş değilse, repository'den bu key ile daha önce kaydedilmiş bir kayıt olup olmadığını sorgula.
		record, err := s.repo.GetIdempotencyRecord(ctx, idempotencyKey)
		if err != nil {
			return err
		}

		// ! "duplicate request"
		if record != nil {
			return nil
		}
	}

	// ! 2. Optimistic Locking Retry Döngüsü
	for {

		wallet, err := s.repo.GetByID(ctx, walletID)
		if err != nil {
			return err
		}

		// * YETKİ KONTROLÜ: Para yatırılacak cüzdan bu kullanıcıya mı ait?
		if wallet.OwnerID != userID {
			return domain.ErrorUnauthorized
		}

		// ! Cüzdan'a para yatır(deposit), hata varsa hatayı dön.
		if err := wallet.Deposit(amount); err != nil {
			return err
		}

		err = s.repo.Update(ctx, wallet)
		if err != nil {
			if errors.Is(err, domain.ErrConcurrentModification) {
				continue
			}
			fmt.Printf("ERROR: Update failed: %v\n", err)
			return err
		}

		// ! Transaction Instance Create and Save
		tn := &domain.Transaction{
			ID:        transactionID,
			WalletID:  walletID,
			Amount:    amount,
			Type:      domain.Deposit,
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveTransaction(ctx, tn); err != nil {
			return err
		}

		break
	}

	// ! 3. Başarılı olan işlemi "Idempotency Kaydı" olarak saklamak.
	if idempotencyKey != "" {
		record := &domain.IdempotencyRecord{
			Key:       idempotencyKey,
			Response:  []byte("Para Yatırma İşlemi Başarılı(SUCCESS)"),
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveIdempotencyRecord(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func (s *walletService) Withdraw(ctx context.Context, idempotencyKey string, walletID string, userID string, transactionID string, amount int64) error {

	if amount <= 0 {
		return domain.ErrorInvalidAmount // Guard Clause
	}

	// ! Idempotency Kontrolü
	if idempotencyKey != "" {
		record, err := s.repo.GetIdempotencyRecord(ctx, idempotencyKey)
		if err != nil {
			return err
		}

		if record != nil {
			return nil
		}
	}

	for {
		// ! Esas İşlemler
		wallet, err := s.repo.GetByID(ctx, walletID)
		if err != nil {
			return err
		}

		// * YETKİ KONTROLÜ
		if wallet.OwnerID != userID {
			return domain.ErrorUnauthorized
		}

		if err := wallet.Withdraw(amount); err != nil {
			return err
		}

		err = s.repo.Update(ctx, wallet)
		if err != nil {
			if errors.Is(err, domain.ErrConcurrentModification) {
				continue
			}
			return err
		}

		// ! Transaction Instance Create and Save
		tn := &domain.Transaction{
			ID:        transactionID,
			WalletID:  walletID,
			Amount:    amount,
			Type:      domain.Withdraw,
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveTransaction(ctx, tn); err != nil {
			return err
		}

		break
	}

	if idempotencyKey != "" {
		record := &domain.IdempotencyRecord{
			Key:       idempotencyKey,
			Response:  []byte("Para Çekme İşlemi Başarılı(SUCCESS)"),
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveIdempotencyRecord(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func (s *walletService) GetTransactions(ctx context.Context, walletID string) ([]*domain.Transaction, error) {

	return s.repo.GetTransactionsByWalletID(ctx, walletID)
}
