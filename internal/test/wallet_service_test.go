package test

import (
	"context"
	"go-hexagonal/internal/adapters/repository"
	"go-hexagonal/internal/core/domain"
	services "go-hexagonal/internal/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
- Deposit Test
- TestWalletService_Deposit = WalletService içerisinde ki Deposit fonksiyonunu test et.
*/
func TestWalletService_Deposit(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	// ? New Wallet Creating..
	wallet, err := service.CreateWallet(
		ctx,
		"Gökhan",
		"TRY",
	)

	// ? "require" → hata varsa testi DURDURUR.
	require.NoError(t, err)

	// ? Wallet.ID numaralı cüzdana 500 yatır.
	err = service.Deposit(
		ctx,
		wallet.ID,
		100,
	)

	require.NoError(t, err)

	// ? Repository’den güncel veri çekmek.
	updated, err := repo.GetByID(
		ctx,
		wallet.ID,
	)

	require.NoError(t, err)

	// ? Okuma sırasında hata var mı?
	assert.Equal(
		t,
		int64(100),
		updated.Balance,
	)
}

/*
- Başarılı Para Çekme Testi
- TestWalletService_Withdraw_Success = Cüzdandan başarılı bir şekilde para çekilmesini test eder.
*/
func TestWalletService_Withdraw_Success(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	// ? Cüzdan oluştur.
	wallet, _ := service.CreateWallet(
		ctx,
		"Can",
		"TRY",
	)

	require.NoError(
		t,
		service.Deposit(
			ctx,
			wallet.ID,
			500,
		),
	)

	err := service.Withdraw(
		ctx,
		wallet.ID,
		200,
	)

	require.NoError(t, err)

	updated, _ := repo.GetByID(
		ctx,
		wallet.ID,
	)

	assert.Equal(
		t,
		int64(300),
		updated.Balance,
	)
}

/*
- Yetersiz Bakiye Testi
- TestWalletService_Withdraw_InsufficientFunds = Cüzdandaki bakiyeden fazla para çekilmeye çalışıldığında yetersiz bakiye hatası alınmasını test eder.
*/
func TestWalletService_Withdraw_InsufficientFunds(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(
		ctx,
		"CAN",
		"TRY",
	)

	err := service.Withdraw(
		ctx,
		wallet.ID,
		300,
	)

	require.ErrorIs(
		t,
		err,
		domain.ErrorInsufficientFunds,
	)
}

/*
- Geçersiz Miktar Para Yatırma Testi
- TestWalletService_Deposit_InvalidAmount = Sıfır veya negatif miktarda para yatırılmaya çalışıldığında geçersiz miktar hatası alınmasını test eder.
*/
func TestWalletService_Deposit_InvalidAmount(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(
		ctx,
		"Gökhan",
		"TRY",
	)

	// "0" yatırma işlemi..
	err := service.Deposit(
		ctx,
		wallet.ID,
		0,
	)

	assert.ErrorIs(
		t,
		err,
		domain.ErrorInvalidAmount,
	)

	// "Negatif" yatırma işlemi..
	err = service.Deposit(
		ctx,
		wallet.ID,
		-100,
	)

	assert.ErrorIs(
		t,
		err,
		domain.ErrorInvalidAmount,
	)
}

/*
- Geçersiz Miktar Para Çekme Testi
- TestWalletService_Withdraw_InvalidAmount = Sıfır veya negatif miktarda para çekilmeye çalışıldığında geçersiz miktar hatası alınmasını test eder.
*/
func TestWalletService_Withdraw_InvalidAmount(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(
		ctx,
		"Gizem",
		"TRY",
	)

	// "0" TL çekme işlemi..
	err := service.Withdraw(
		ctx,
		wallet.ID,
		0,
	)

	assert.ErrorIs(
		t,
		err,
		domain.ErrorInvalidAmount,
	)

	err = service.Withdraw(
		ctx,
		wallet.ID,
		-50,
	)

	assert.ErrorIs(
		t,
		err,
		domain.ErrorInvalidAmount,
	)
}
