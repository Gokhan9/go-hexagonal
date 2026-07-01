package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-hexagonal/internal/core/domain"
)

type PostgreWalletRepository struct {
	db *sql.DB
}

// NewPostgresWalletRepository veritabanı bağlantısıyla yeni bir repo instance'ı döner
func NewPostgreWalletRepository(db *sql.DB) *PostgreWalletRepository {
	return &PostgreWalletRepository{
		db: db,
	}
}

// CREATE, Yeni bir "wallet" kaydeder
func (r *PostgreWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {

	query := `
	INSERT INTO wallets (id, owner_id, balance, currency, version, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.getExecutor(ctx).ExecContext(ctx, query,
		wallet.ID,
		wallet.OwnerID,
		wallet.Balance,
		wallet.Currency,
		wallet.Version,
		wallet.CreatedAt,
	)
	return err
}

// GetByID Id'ye göre wallet getirir.
func (r *PostgreWalletRepository) GetByID(ctx context.Context, id string) (*domain.Wallet, error) {

	query := `SELECT id, owner_id, balance, currency, version, created_at FROM wallets WHERE id= $1`

	row := r.getExecutor(ctx).QueryRowContext(ctx, query, id)

	var wallet domain.Wallet
	err := row.Scan(
		&wallet.ID,
		&wallet.OwnerID,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.Version,
		&wallet.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrorWalletNotFound
	}

	if err != nil {
		return nil, err
	}

	return &wallet, nil

}

func (r *PostgreWalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {

	query := `UPDATE wallets SET balance = $1, version = version + 1 WHERE id = $2 AND version = $3`

	// result, err := r.db.ExecContext(ctx, query, wallet.Balance, wallet.ID, wallet.Version)
	result, err := r.getExecutor(ctx).ExecContext(ctx, query, wallet.Balance, wallet.ID, wallet.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrConcurrentModification
	}

	// 200 ise domain nesnesini update et
	wallet.Version++
	return nil
}

func (r *PostgreWalletRepository) getExecutor(ctx context.Context) DBExecutor {

	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *PostgreWalletRepository) BeginTx(ctx context.Context) (context.Context, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return WithTx(ctx, tx), nil // context içine tx'i gömüyoruz
}

func (r *PostgreWalletRepository) Commit(ctx context.Context) error {
	tx := GetTx(ctx)
	if tx == nil {
		return errors.New("no transaction found in context.")
	}

	return tx.Commit()
}

func (r *PostgreWalletRepository) Rollback(ctx context.Context) error {
	tx := GetTx(ctx)
	if tx == nil {
		return nil
	}
	return tx.Rollback()
}
