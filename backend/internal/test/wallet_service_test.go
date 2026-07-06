package test

import (
	"context"
	"go-hexagonal/internal/adapters/repository"
	"go-hexagonal/internal/core/domain"
	"go-hexagonal/internal/core/service"
	services "go-hexagonal/internal/core/service"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletService_Deposit(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, err := service.CreateWallet(ctx, "Gökhan", "TRY")
	require.NoError(t, err)

	err = service.Deposit(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 100)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(100), updated.Balance)
}

func TestWalletService_Withdraw_Success(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(ctx, "Can", "TRY")

	require.NoError(t, service.Deposit(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 500))

	err := service.Withdraw(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 200)
	require.NoError(t, err)

	updated, _ := repo.GetByID(ctx, wallet.ID)
	assert.Equal(t, int64(300), updated.Balance)
}

func TestWalletService_Withdraw_InsufficientFunds(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(ctx, "CAN", "TRY")

	err := service.Withdraw(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 300)
	require.ErrorIs(t, err, domain.ErrorInsufficientFunds)
}

func TestWalletService_Deposit_InvalidAmount(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(ctx, "Gökhan", "TRY")

	err := service.Deposit(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 0)
	assert.ErrorIs(t, err, domain.ErrorInvalidAmount)

	err = service.Deposit(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), -100)
	assert.ErrorIs(t, err, domain.ErrorInvalidAmount)
}

func TestWalletService_Withdraw_InvalidAmount(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(ctx, "Gizem", "TRY")

	err := service.Withdraw(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 0)
	assert.ErrorIs(t, err, domain.ErrorInvalidAmount)

	err = service.Withdraw(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), -50)
	assert.ErrorIs(t, err, domain.ErrorInvalidAmount)
}

func TestWalletService_Deposit_Concurrent(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(ctx, "Gökhan", "TRY")

	const goroutineCount = 100
	const depositAmount = 10

	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()
			_ = service.Deposit(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), depositAmount)
		}()
	}

	wg.Wait()

	updated, _ := repo.GetByID(ctx, wallet.ID)
	assert.Equal(t, int64(goroutineCount*depositAmount), updated.Balance)
}

func TestWalletService_Deposit_Idempotency(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, err := service.CreateWallet(ctx, "Hakan", "TRY")
	require.NoError(t, err)

	idempotencyKey := "unique-deposit-key-1"

	err = service.Deposit(ctx, idempotencyKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 1000)
	require.NoError(t, err)

	err = service.Deposit(ctx, idempotencyKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 1000)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1000), updated.Balance)
}

func TestWalletService_Withdraw_Idempotency(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, err := service.CreateWallet(ctx, "Mert", "TRY")
	require.NoError(t, err)

	err = service.Deposit(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 2000)
	require.NoError(t, err)

	idempotencyKey := "unique-withdraw-key-1"

	err = service.Withdraw(ctx, idempotencyKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 500)
	require.NoError(t, err)

	err = service.Withdraw(ctx, idempotencyKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 500)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1500), updated.Balance)
}

func TestWalletService_TransactionHistory_Verification(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, err := service.CreateWallet(ctx, "Gökhan", "TRY")
	require.NoError(t, err)

	err = service.Deposit(ctx, "key-verify-deposit", wallet.ID, wallet.OwnerID, uuid.NewString(), 1000)
	require.NoError(t, err)

	err = service.Withdraw(ctx, "key-verify-withdraw", wallet.ID, wallet.OwnerID, uuid.NewString(), 300)
	require.NoError(t, err)

	tns, err := repo.GetTransactionsByWalletID(ctx, wallet.ID)
	require.NoError(t, err)

	assert.Len(t, tns, 2)
	assert.Equal(t, domain.Deposit, tns[0].Type)
	assert.Equal(t, int64(1000), tns[0].Amount)
	assert.Equal(t, domain.Withdraw, tns[1].Type)
	assert.Equal(t, int64(300), tns[1].Amount)
}

func TestWalletService_Full_E2E_Scenario(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, err := service.CreateWallet(ctx, "Gökhan", "TRY")
	require.NoError(t, err)
	assert.Equal(t, int64(0), wallet.Balance)

	depositKey := "unique-deposit-key-999"
	err = service.Deposit(ctx, depositKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 10050)
	require.NoError(t, err)

	err = service.Deposit(ctx, depositKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 10050)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(10050), updated.Balance)

	withdrawKey := "unique-withdraw-key-999"
	err = service.Withdraw(ctx, withdrawKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 3050)
	require.NoError(t, err)

	err = service.Withdraw(ctx, withdrawKey, wallet.ID, wallet.OwnerID, uuid.NewString(), 3050)
	require.NoError(t, err)

	updated, err = repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(7000), updated.Balance)

	tns, err := service.GetTransactions(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Len(t, tns, 2)
	assert.Equal(t, domain.Deposit, tns[0].Type)
	assert.Equal(t, int64(10050), tns[0].Amount)
	assert.Equal(t, domain.Withdraw, tns[1].Type)
	assert.Equal(t, int64(3050), tns[1].Amount)
}

func TestWalletService_Transfer_Success(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := service.NewWalletService(repo)
	ctx := context.Background()

	// Wallet Create
	w1, _ := service.CreateWallet(ctx, "Gökhan", "TRY")
	w2, _ := service.CreateWallet(ctx, "Can", "TRY")

	// Bakiye yükle
	service.Deposit(ctx, "", w1.ID, w1.OwnerID, uuid.NewString(), 1000)
	service.Deposit(ctx, "", w2.ID, w2.OwnerID, uuid.NewString(), 500)

	// Transfer
	err := service.Transfer(ctx, "transfer-key-1", w1.ID, w2.ID, w1.OwnerID, 300)
	require.NoError(t, err)

	// Doğrulama
	updatedW1, _ := repo.GetByID(ctx, w1.ID)
	updatedW2, _ := repo.GetByID(ctx, w2.ID)

	assert.Equal(t, int64(700), updatedW1.Balance)
	assert.Equal(t, int64(800), updatedW2.Balance)
}

func TestWalletService_Transfer_InsufficientFunds(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := service.NewWalletService(repo)
	ctx := context.Background()

	w1, _ := service.CreateWallet(ctx, "Gökhan", "TRY")
	w2, _ := service.CreateWallet(ctx, "Can", "TRY")

	// 100 try var
	service.Deposit(ctx, "", w1.ID, w1.OwnerID, uuid.NewString(), 100)

	// 200 transfer et
	err := service.Transfer(ctx, "transfer-key-2", w1.ID, w2.ID, w1.OwnerID, 200)
	require.ErrorIs(t, err, domain.ErrorInsufficientFunds)

	// bakiye aynı kalmalı
	updatedW1, _ := repo.GetByID(ctx, w1.ID)
	updatedW2, _ := repo.GetByID(ctx, w2.ID)

	assert.Equal(t, int64(100), updatedW1.Balance)
	assert.Equal(t, int64(0), updatedW2.Balance)
}

func TestWalletService_GetBalance(t *testing.T) {
	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	// 1. Senaryo: create wallet
	wallet, err := service.CreateWallet(ctx, "Gökhan", "TRY")
	require.NoError(t, err)

	// bakiye ekle(250)
	service.Deposit(ctx, "", wallet.ID, wallet.OwnerID, uuid.NewString(), 250)

	balance, err := service.GetBalance(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(250), balance)

	// 2. Senaryo: Olmayan cüzdan için hata kontrolü
	_, err = service.GetBalance(ctx, "non-existend-id")
	require.EqualError(t, err, "wallet not found..")
}
