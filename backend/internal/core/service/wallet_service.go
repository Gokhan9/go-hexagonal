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

func (s *walletService) Deposit(ctx context.Context, idempotencyKey string, walletID string, userID string, transactionID string, amount int64) error {

	txContext, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer s.repo.Rollback(txContext)

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

	return s.repo.Commit(txContext)
}

func (s *walletService) Withdraw(ctx context.Context, idempotencyKey string, walletID string, userID string, transactionID string, amount int64) error {

	if amount <= 0 {
		return domain.ErrorInvalidAmount // Guard Clause
	}

	// Transaction Start
	txContext, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer s.repo.Rollback(txContext)

	// Idempotency Kontrolü
	if idempotencyKey != "" {
		record, err := s.repo.GetIdempotencyRecord(txContext, idempotencyKey)
		if err != nil {
			return err
		}

		if record != nil {
			return nil
		}
	}

	// 2. Retry Döngüsü
	for {
		wallet, err := s.repo.GetByID(txContext, walletID)
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

		// Update (Concurrency Check)
		err = s.repo.Update(txContext, wallet)
		if err != nil {
			if errors.Is(err, domain.ErrConcurrentModification) {
				continue
			}
			return err
		}

		// Transaction Instance Create and Save
		tn := &domain.Transaction{
			ID:        transactionID,
			WalletID:  walletID,
			Amount:    amount,
			Type:      domain.Withdraw,
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveTransaction(txContext, tn); err != nil {
			return err
		}

		break
	}

	// Idempotency Kaydı (Success Durumunda)
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

	// 3. Commit
	return s.repo.Commit(txContext)
}

func (s *walletService) GetTransactions(ctx context.Context, walletID string) ([]*domain.Transaction, error) {

	return s.repo.GetTransactionsByWalletID(ctx, walletID)
}

func (s *walletService) Transfer(ctx context.Context, idempotencyKey, fromWalletID, toWalletID, ownerID string, amount int64) error {

	// guard clause
	if fromWalletID == toWalletID {
		return domain.ErrorSelfTransfer
	}

	if amount <= 0 {
		return domain.ErrorInvalidAmount
	}

	// Transaction Start
	// NOT: BeginTx, transaction'ı ctx içerisine gömer (context-based propagation)
	ctx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Transaction tamamlanmazsa "panic&hata" rollback'i dön. Database tutarlılığını sağlar.
	defer s.repo.Rollback(ctx)

	// * Business Logics

	// Gönderen(withdraw) wallet'ı çek
	fromWallet, err := s.repo.GetByID(ctx, fromWalletID)
	if err != nil {
		return err
	}

	// Gönderen(withdraw) wallet'ın sahibi mi kontrolü
	if fromWallet.OwnerID != ownerID {
		return domain.ErrorUnauthorized
	}

	// Alıcı(deposit) wallet'ı çek.
	toWallet, err := s.repo.GetByID(ctx, toWalletID)
	if err != nil {
		return err
	}

	// balance check
	if fromWallet.Balance < amount {
		return domain.ErrorInsufficientFunds
	}

	// İşlemleri Gerçekleştir.

	// withdraw(gönderen)
	fromWallet.Balance -= amount
	if err := s.repo.Update(ctx, fromWallet); err != nil {
		return err
	}

	// deposit(alıcı)
	toWallet.Balance += amount
	if err := s.repo.Update(ctx, toWallet); err != nil {
		return err
	}

	// commit (işlemlerin başarılı olması durumunda.)
	return s.repo.Commit(ctx)
}

// GetByID ile wallet'a ait balance(bakiye) bilgisini döner.
func (s *walletService) GetBalance(ctx context.Context, walletID string) (int64, error) {

	wallet, err := s.repo.GetByID(ctx, walletID)

	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}
