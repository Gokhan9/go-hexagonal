package test

import (
	"context"
	"go-hexagonal/internal/adapters/repository"
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
