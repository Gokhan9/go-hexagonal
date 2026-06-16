package services

import (
	"context"
	"errors"
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

// ! 2. ADIM : WalletService, "Deposit ve Withdraw" Güncellemesi
func (s *walletService) Deposit(ctx context.Context, idempotencyKey string, walletID string, amount int64) error {

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

		// ! Cüzdan'a para yatır(deposit), hata varsa hatayı dön.
		if err := wallet.Deposit(amount); err != nil {
			return err
		}

		err = s.repo.Update(ctx, wallet)
		if err != nil {
			if errors.Is(err, domain.ErrConcurrentModification) {
				continue
			}

			break // ! Update Başarılı, döngüden çıkar.
		}

		// Eşzamanlılık(Concurrency) hatası alındıysa döngü başa döner ve tekrar dener. // Güncel Cüzdanı (ve yeni versiyonunu) tekrar çekip yeniden dener.

		/*
			Deposit isteği
				 ↓
			Transaction nesnesi oluştur
				 ↓
			Transaction.ID = rastgele UUID üret
				 ↓
			SaveTransaction()

		*/
		// ! Transaction Instance Create and Save - 16.06.2026
		tn := &domain.Transaction{
			ID:        uuid.NewString(), // "Transaction Kaydına" benzersiz(unique) kimlik (ID) vermek için kullanırız. (örn:"d6d0b8b8-76ab-4f7a-b56c-8d3d0c11c4df")
			WalletID:  walletID,
			Amount:    amount,
			Type:      domain.Deposit,
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveTransaction(ctx, tn); err != nil {
			return err // ! İşlem Kaydı(Transaction), başarısızsa akışı kesiyoruz.
		}
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

func (s *walletService) Withdraw(ctx context.Context, idempotencyKey string, walletID string, amount int64) error {

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

		if err := wallet.Withdraw(amount); err != nil {
			return err
		}

		// "Transaction" Create and Save işlemi "unreachable code" uyarısı alıyorum, alt satırı == çevirince kod akışı düzeldi.
		err = s.repo.Update(ctx, wallet)
		if err == nil {
			break
		}

		tn := &domain.Transaction{
			ID:        uuid.NewString(),
			WalletID:  walletID,
			Amount:    amount,
			Type:      domain.Withdraw,
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveTransaction(ctx, tn); err != nil {
			return err
		}
	}

	// ! 3. Başarılı İşlem - Idempotency Kaydı Olarak Sakla
	if idempotencyKey != "" {
		record := &domain.IdempotencyRecord{
			Key:       idempotencyKey,
			Response:  []byte("Para Çekme İşlemi Başarılı(SUCCESS)"),
			CreatedAt: time.Now(),
		}

		if err := s.repo.SaveIdempotencyRecord(ctx, record); err != nil {
			return nil
		}
	}

	return nil
}

func (s *walletService) GetTransactions(ctx context.Context, walletID string) ([]*domain.Transaction, error) {

	return s.repo.GetTransactionsByWalletID(ctx, walletID)
}
