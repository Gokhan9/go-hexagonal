package test

import (
	"context"
	"database/sql"
	"fmt"
	"go-hexagonal/internal/adapters/repository"
	"go-hexagonal/internal/core/service"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_FinancialWorkflow(t *testing.T) {
	// 1. Setup: Gerçek veritabanı bağlantısı
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=wallet_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Tabloları temizle (Sıfırdan test için)
	_, err = db.Exec("TRUNCATE TABLE transactions, audit_logs, wallets, idempotency_records RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	walletRepo := repository.NewPostgreWalletRepository(db)
	auditRepo := repository.NewPostgreAuditRepository(db)
	svc := service.NewWalletService(walletRepo, auditRepo)
	ctx := context.Background()

	// 2. Senaryo: İki Cüzdan Oluştur
	w1, err := svc.CreateWallet(ctx, "Gökhan", "TRY")
	require.NoError(t, err)
	w2, err := svc.CreateWallet(ctx, "Hakan", "TRY")
	require.NoError(t, err)

	// 3. Deposit (Gökhan)
	txID1 := uuid.NewString()
	fmt.Printf("DEBUG: Depositing with txID: %s\n", txID1) // DEBUG
	err = svc.Deposit(ctx, "", w1.ID, "Gökhan", txID1, 1000)
	require.NoError(t, err)

	// 4. Transfer (Gökhan -> Hakan)
	err = svc.Transfer(ctx, "", w1.ID, w2.ID, "Gökhan", 400)
	require.NoError(t, err)

	// 5. Doğrulama
	bal1, _ := svc.GetBalance(ctx, w1.ID)
	bal2, _ := svc.GetBalance(ctx, w2.ID)

	assert.Equal(t, int64(600), bal1)
	assert.Equal(t, int64(400), bal2)

	// Transaction Status Kontrolü
	var status string
	err = db.QueryRow("SELECT status FROM transactions WHERE id = $1", txID1).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "COMPLETED", status)

	// Audit Log Kontrolü
	var op string
	err = db.QueryRow("SELECT operation FROM audit_logs WHERE entity_id = $1", w1.ID).Scan(&op)
	require.NoError(t, err)
	assert.Equal(t, "TRANSFER", op)

	fmt.Println("E2E Smoke Test Başarıyla Tamamlandı!")
}
