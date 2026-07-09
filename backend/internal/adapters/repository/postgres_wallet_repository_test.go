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
