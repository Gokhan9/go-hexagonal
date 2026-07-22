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
	walletRepo ports.WalletRepository
	auditRepo  ports.AuditRepository // NEW
}

func NewWalletService(repo ports.WalletRepository, auditRepo ports.AuditRepository) ports.WalletService {
	return &walletService{
		walletRepo: repo,
		auditRepo:  auditRepo, // NEW
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

	if err := s.walletRepo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) GetWallet(ctx context.Context, userID, id string) (*domain.Wallet, error) {

	wallet, err := s.walletRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// * 2- YETKI KONTROLÜ: Bu wallet istek atan user'a mı ait?
	//if wallet.OwnerID != userID {
	//	return nil, domain.ErrorUnauthorized
	//}

	return wallet, nil
}

func (s *walletService) Deposit(ctx context.Context, idempotencyKey string, walletID string, userID string, transactionID string, amount int64) error {

	txContext, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return err
	}

	// ! Hata durumunda işlemi FAILED olarak işaretleyip rollback işlemi yapmak.
	defer func() {
		if err != nil {
			_ = s.walletRepo.UpdateTransactionStatus(txContext, transactionID, domain.StatusFailed)
			_ = s.walletRepo.Rollback(txContext)
		}
	}()

	if amount <= 0 {
		return domain.ErrorInvalidAmount // Guard Clause
	}

	tn := &domain.Transaction{
		ID:        transactionID,
		WalletID:  walletID,
		Amount:    amount,
		Type:      domain.Deposit,
		Status:    domain.StatusPending, // "Pending", başlangıç
		CreatedAt: time.Now(),
	}

	if err = s.walletRepo.SaveTransaction(txContext, tn); err != nil {
		return err
	}

	// 1. Idempotency Kontrolü
	if idempotencyKey != "" {
		record, err := s.walletRepo.GetIdempotencyRecord(txContext, idempotencyKey)
		if err != nil {
			return err
		}

		// "duplicate request"
		if record != nil {
			return nil
		}
	}

	// 2. Optimistic Locking Retry Döngüsü
	for {
		wallet, err := s.walletRepo.GetByID(txContext, walletID)
		if err != nil {
			return err
		}

		fmt.Printf("DEBUG FOR Deposit: WalletID: %s, WalletOwnerID: '%s', RequestUserID: '%s'\n", walletID, wallet.OwnerID, userID)

		// * YETKİ KONTROLÜ
		if wallet.OwnerID != userID {
			return domain.ErrorUnauthorized
		}

		if err := wallet.Deposit(amount); err != nil {
			return err
		}

		err = s.walletRepo.Update(txContext, wallet)
		if err != nil {
			if errors.Is(err, domain.ErrConcurrentModification) {
				continue
			}
			fmt.Printf("ERROR: Update failed: %v\n", err)
			return err
		}
		break
	}

	// ! İşlem biter, Status = "COMPLETED" olarak güncellenir.
	if err = s.walletRepo.UpdateTransactionStatus(txContext, transactionID, domain.StatusCompleted); err != nil {
		return err
	}

	// 3. Başarılı olan işlemi "Idempotency Kaydı" olarak saklamak.
	if idempotencyKey != "" {
		record := &domain.IdempotencyRecord{
			Key:       idempotencyKey,
			Response:  []byte("Para Yatırma İşlemi Başarılı(SUCCESS)"),
			CreatedAt: time.Now(),
		}

		if err := s.walletRepo.SaveIdempotencyRecord(txContext, record); err != nil {
			return err
		}
	}

	return s.walletRepo.Commit(txContext)
}

func (s *walletService) Withdraw(ctx context.Context, idempotencyKey string, walletID string, userID string, transactionID string, amount int64) (err error) {

	// Transaction Start
	txContext, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = s.walletRepo.UpdateTransactionStatus(txContext, transactionID, domain.StatusFailed)
			_ = s.walletRepo.Rollback(txContext)
		}
	}()

	if amount <= 0 {
		return domain.ErrorInvalidAmount // Guard Clause
	}

	tn := &domain.Transaction{
		ID:        transactionID,
		WalletID:  walletID,
		Amount:    amount,
		Type:      domain.Withdraw,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	if err = s.walletRepo.SaveTransaction(txContext, tn); err != nil {
		return err
	}

	// Idempotency Kontrolü
	if idempotencyKey != "" {
		record, err := s.walletRepo.GetIdempotencyRecord(txContext, idempotencyKey)
		if err != nil {
			return err
		}

		if record != nil {
			return nil
		}
	}

	// 2. Retry Döngüsü
	for {
		wallet, err := s.walletRepo.GetByID(txContext, walletID)
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
		err = s.walletRepo.Update(txContext, wallet)
		if err != nil {
			if errors.Is(err, domain.ErrConcurrentModification) {
				continue
			}
			return err
		}
		break
	}

	if err = s.walletRepo.UpdateTransactionStatus(txContext, transactionID, domain.StatusCompleted); err != nil {
		return err
	}

	// Idempotency Kaydı (Success Durumunda)
	if idempotencyKey != "" {
		record := &domain.IdempotencyRecord{
			Key:       idempotencyKey,
			Response:  []byte("Para Çekme İşlemi Başarılı(SUCCESS)"),
			CreatedAt: time.Now(),
		}

		if err := s.walletRepo.SaveIdempotencyRecord(txContext, record); err != nil {
			return err
		}
	}

	// 3. Commit
	return s.walletRepo.Commit(txContext)
}

func (s *walletService) GetTransactions(ctx context.Context, walletID string) ([]*domain.Transaction, error) {

	return s.walletRepo.GetTransactionsByWalletID(ctx, walletID)
}

func (s *walletService) Transfer(ctx context.Context, idempotencyKey, fromWalletID, toWalletID, ownerID string, amount int64) (err error) {

	// guard clause
	if fromWalletID == toWalletID {
		return domain.ErrorSelfTransfer
	}

	if amount <= 0 {
		return domain.ErrorInvalidAmount
	}

	// Transaction Start
	txContext, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return err
	}

	transactionID := uuid.NewString() // Transfer için yeni ID
	defer func() {
		if err != nil {
			_ = s.walletRepo.UpdateTransactionStatus(txContext, transactionID, domain.StatusFailed)
			_ = s.walletRepo.Rollback(txContext)
		}
	}()

	tn := &domain.Transaction{
		ID:        transactionID,
		WalletID:  fromWalletID,
		Amount:    amount,
		Type:      domain.Withdraw,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}
	if err = s.walletRepo.SaveTransaction(txContext, tn); err != nil {
		return err
	}

	// * BUSINESS LOGICS

	// Gönderen(withdraw) wallet'ı çek
	fromWallet, err := s.walletRepo.GetByID(txContext, fromWalletID)
	if err != nil {
		return err
	}

	// Gönderen(withdraw) wallet'ın sahibi mi kontrolü
	if fromWallet.OwnerID != ownerID {
		return domain.ErrorUnauthorized
	}

	// Alıcı(deposit) wallet'ı çek.
	toWallet, err := s.walletRepo.GetByID(txContext, toWalletID)
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
	if err := s.walletRepo.Update(txContext, fromWallet); err != nil {
		return err
	}

	// deposit(alıcı)
	toWallet.Balance += amount
	if err := s.walletRepo.Update(txContext, toWallet); err != nil {
		return err
	}

	// Başarılı: Statüyü COMPLETED olarak güncelle
	if err = s.walletRepo.UpdateTransactionStatus(txContext, transactionID, domain.StatusCompleted); err != nil {
		return err
	}

	// COMMIT (İşlemlerin başarılı olması durumunda.)
	err = s.walletRepo.Commit(txContext)
	if err != nil {
		return err
	}

	// ! AUDIT LOGLAMA
	s.auditRepo.Save(ctx, &domain.AuditLog{
		ID:         uuid.NewString(),
		EntityID:   fromWalletID,
		EntityType: "WALLET",
		Operation:  "TRANSFER",
		UserID:     ownerID,
		Changes:    []byte(fmt.Sprintf("Transferred %d from %s to %s", amount, fromWalletID, toWalletID)),
		CreatedAt:  time.Now(),
	})

	return nil
}

// GetByID ile wallet'a ait balance(bakiye) bilgisini döner.
func (s *walletService) GetBalance(ctx context.Context, walletID string) (int64, error) {

	wallet, err := s.walletRepo.GetByID(ctx, walletID)

	if err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

func (s *walletService) CloseWallet(ctx context.Context, walletID string, userID string) error {

	// Optimistic Lock Retry Döngüsü
	for {
		// 1. Wallet'ı çek
		wallet, err := s.walletRepo.GetByID(ctx, walletID)
		if err != nil {
			return err
		}

		// 2. wallet için yetki kontrolü
		if wallet.OwnerID != userID {
			return domain.ErrorUnauthorized
		}

		// 3. wallet için logicler
		if wallet.Status == domain.StatusClosed {
			return domain.ErrorWalletAlreadyClosed
		}
		if wallet.Balance != 0 {
			return domain.ErrorWalletNotEmptied
		}

		// 4. Optimistic lock ile durumu güncelle
		err = s.walletRepo.UpdateWalletStatus(ctx, walletID, domain.StatusClosed, wallet.Version)
		if err != nil {
			// Concurrency hatası alsak dahi döngüye devam ederiz.
			if errors.Is(err, domain.ErrConcurrentModification) {
				continue
			}

			return err
		}
		// Success
		break
	}
	return nil

}
