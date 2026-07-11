package repository

import (
	"context"
	"go-hexagonal/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresWalletRepository_CreateAndGet(t *testing.T) {

	// 1. Setup(Container başlat ve db'yi hazırlamak)
	db, cleanup := setupDatabase(t)
	defer cleanup() // Test biter, container kapanır.

	repo := NewPostgreWalletRepository(db)
	ctx := context.Background()

	// 2. Test Verisi
	wallet := &domain.Wallet{
		ID:        uuid.NewString(),
		OwnerID:   "user1",
		Balance:   100,
		Currency:  "TRY",
		Version:   1,
		CreatedAt: time.Now(),
	}

	// 3. Eylem (CREATE)
	err := repo.Create(ctx, wallet)
	require.NoError(t, err)

	// 4. Doğrulama (GetByID)
	fetch, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)

	assert.Equal(t, wallet.ID, fetch.ID)
	assert.Equal(t, wallet.Balance, fetch.Balance)
	assert.Equal(t, wallet.OwnerID, fetch.OwnerID)
}

func TestPostgresWalletRepository_Transaction(t *testing.T) {

	db, cleanup := setupDatabase(t)
	defer cleanup()

	repo := NewPostgreWalletRepository(db)
	ctx := context.Background()

	wallet := &domain.Wallet{
		ID:        uuid.NewString(),
		OwnerID:   "user_tx",
		Balance:   100,
		Currency:  "TRY",
		Version:   1,
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, wallet)
	require.NoError(t, err)

	t.Run("Commit_Success", func(t *testing.T) {

		txCtx, err := repo.BeginTx(ctx)
		require.NoError(t, err)

		wallet.Balance = 200
		err = repo.Update(txCtx, wallet)
		require.NoError(t, err)

		//commit
		err = repo.Commit(txCtx)
		require.NoError(t, err)

		// Doğrulama: Commit Sonrası Bakiye Güncellenmeli
		updated, err := repo.GetByID(ctx, wallet.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(200), updated.Balance)
	})

	t.Run("Rollback_Success", func(t *testing.T) {

		txCtx, err := repo.BeginTx(ctx)
		require.NoError(t, err)

		wallet.Balance = 500
		err = repo.Update(txCtx, wallet)
		require.NoError(t, err)

		//rollback
		err = repo.Rollback(txCtx)
		require.NoError(t, err)

		current, err := repo.GetByID(ctx, wallet.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(200), current.Balance)
	})
}
