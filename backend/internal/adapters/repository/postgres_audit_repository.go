package repository

import (
	"context"
	"database/sql"
	"go-hexagonal/internal/core/domain"
)

/*
! "postgres_audit_repository", AuditRepository portunun veritabanı üzerindeki karşılığıdır ve tıpkı WalletRepository gibi çalışır.
*/

type PostgreAuditRepository struct {
	db *sql.DB
}

func NewPostgreAuditRepository(db *sql.DB) *PostgreAuditRepository {
	return &PostgreAuditRepository{
		db: db,
	}
}

func (r *PostgreAuditRepository) Save(ctx context.Context, log *domain.AuditLog) error {

	query := ` INSERT INTO audit_logs (id, entity_id, entity_type, operation, user_id, changes, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// Eğer transaction içerisinde çalışıyorsa "tx" kullanılabilir. Bunu wallet_repository'deki "getExecutor" mantığıyla yapabiliriz
	executor := r.getExecutor(ctx)

	_, err := executor.ExecContext(ctx, query,
		log.ID,
		log.EntityID,
		log.EntityType,
		log.Operation,
		log.UserID,
		log.Changes,
		log.CreatedAt,
	)
	return err
}

func (r *PostgreAuditRepository) getExecutor(ctx context.Context) DBExecutor {

	if tx := GetTx(ctx); tx != nil {
		return tx
	}

	return r.db
}
