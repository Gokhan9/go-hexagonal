package test

import (
	"context"
	"go-hexagonal/internal/adapters/repository"
	"go-hexagonal/internal/core/domain"
	services "go-hexagonal/internal/core/service"
	"sync"
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
		"",
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
			"",
			wallet.ID,
			500,
		),
	)

	err := service.Withdraw(
		ctx,
		"",
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
		"",
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
		"",
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
		"",
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
		"",
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
		"",
		wallet.ID,
		-50,
	)

	assert.ErrorIs(
		t,
		err,
		domain.ErrorInvalidAmount,
	)
}

/*
Eşzamanlılık(Concurrency) hatasını yakalamak ve kodun doğruluğunu kanıtlamak için Race Condition Testi yazıyoruz. Aynı anda 100 adet goroutine ile cüzdana para yatırabiliriz.
*/

func TestWalletService_Deposit_Concurrent(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	wallet, _ := service.CreateWallet(ctx, "Gökhan", "TRY")

	const goroutineCount = 100
	const depositAmount = 10 // Her defasında 10 birim yatır..

	var wg sync.WaitGroup
	wg.Add(goroutineCount) // Beklenecek goroutine sayısı(100).

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done() // goroutine işini tamamladığında sayaçtan "1" eksilir. Function başına yazılır hata payını düşürmek için..
			_ = service.Deposit(ctx, "", wallet.ID, depositAmount)
		}()
	}

	wg.Wait() // Sayaç 0 olana kadar diğer işlemleri bloklarız, 0 olduğunda program kaldığı yerden devam edebilir.

	updated, _ := repo.GetByID(ctx, wallet.ID) // Güncel bakiye kontrolü

	assert.Equal(
		t,
		int64(goroutineCount*depositAmount),
		updated.Balance)
}

/*
- Senaryo A (Deposit Idempotency): Aynı idempotency key ile ardışık iki deposit yapıldığında, bakiye yalnızca bir kere artmalı.
*/
func TestWalletService_Deposit_Idempotency(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	// ! Wallet Create
	wallet, err := service.CreateWallet(ctx, "Hakan", "TRY")
	require.NoError(t, err)

	// ! IdempotencyKey
	idempotencyKey := "unique-deposit-key-1"

	// ! İlk Para Yatırma(Deposit) İşlemi Başarılı Olmalı
	err = service.Deposit(ctx, idempotencyKey, wallet.ID, 1000) // 10 TL (1000 kuruş)
	require.NoError(t, err)

	// ! İkinci Para Yatırma(Deposit) İşlemi "IDEMPOTENT" Olmalı - Tekrar İşlem Yapmamalı.
	err = service.Deposit(ctx, idempotencyKey, wallet.ID, 1000)
	require.NoError(t, err) // Hata vermemeli, başarılı gibi "nil" dönmeli.

	// ! Bakiye Kontrolü (Bakiye 20 TL değil, 10 TL olmalı)
	updated, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(
		t,
		int64(1000),
		updated.Balance)
}

/*
- Senaryo B (Withdraw Idempotency): Aynı idempotency key ile ardışık iki withdraw yapıldığında, bakiye yalnızca bir kere azalmalı.
*/
func TestWalletService_Withdraw_Idempotency(t *testing.T) {

	repo := repository.NewMemoryWalletRepository()
	service := services.NewWalletService(repo)
	ctx := context.Background()

	// ! Wallet Create
	wallet, err := service.CreateWallet(ctx, "Mert", "TRY")
	require.NoError(t, err)

	// ! İlk olarak bakiye yükleyelim (idempotency key kullanmadan geçebiliriz)
	err = service.Deposit(ctx, "", wallet.ID, 2000)
	require.NoError(t, err)

	idempotencyKey := "unique-deposit-key-1"

	// ! İkinci olarak İlk Para Çekme İşlemi (Başarılı Olmalı)
	err = service.Withdraw(ctx, idempotencyKey, wallet.ID, 500)
	require.NoError(t, err)

	// !  İkinci Para Çekme İşlemi - (Idempotent Olmalı - Tekrar para çekmemeli)
	err = service.Withdraw(ctx, idempotencyKey, wallet.ID, 500)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	assert.Equal(
		t,
		int64(1500),
		updated.Balance,
	)
}
