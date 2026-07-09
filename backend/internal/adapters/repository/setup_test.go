package repository

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupDatabase(t *testing.T) (*sql.DB, func()) {

	ctx := context.Background()

	// 1. POSTGRESQL Container Start
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connection").
				WithOccurrence(2)),
	)

	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// 2. DB BAĞLANTISI KURMA
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	// 3. MIGRATION ÇALIŞTIR
	migrationPath := filepath.Join("..", "..", "..", "migrations", "001_init_schema.sql")
	schema, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("failed to read migration file: %v", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	// cleanup function: TEST BİTİNCE KONTEYNERI DURDURUR.
	cleanup := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}
