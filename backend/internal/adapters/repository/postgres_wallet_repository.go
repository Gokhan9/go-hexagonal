package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-hexagonal/internal/core/domain"
	"log"
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

	query := `INSERT INTO wallets (id, owner_id, balance, currency, version, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

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

/*
1. "GetByID", verilen wallet ID'sine ait cüzdan kaydını veritabanından getirir.
2. Context'te transaction varsa sorguyu transaction üzerinden yürütür.
3. Gelen satırı domain.Wallet struct'ına map eder.
4. Kayıt yoksa ErrorWalletNotFound döndürür.
5. Başarılıysa Wallet nesnesini döndürür.
*/
func (r *PostgreWalletRepository) GetByID(ctx context.Context, id string) (*domain.Wallet, error) {

	// 1
	query := `SELECT id, owner_id, balance, currency, version, created_at FROM wallets WHERE id= $1`

	// 2.
	row := r.getExecutor(ctx).QueryRowContext(ctx, query, id)

	// 3.
	var wallet domain.Wallet
	err := row.Scan(
		&wallet.ID,
		&wallet.OwnerID,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.Version,
		&wallet.CreatedAt,
	)

	// 4.
	if err == sql.ErrNoRows {
		return nil, domain.ErrorWalletNotFound
	}

	if err != nil {
		return nil, err
	}

	// 5.
	return &wallet, nil

}

/*
1. "Update", wallet'ın bakiyesini optimistic locking kullanarak günceller.
2. Transaction varsa transaction üzerinden, yoksa normal DB bağlantısı üzerinden sorguyu çalıştır.
3. Eğer version uyuşmazsa (başka bir işlem kaydı değiştirmişse) domain.ErrConcurrentModification hatası döndürülür.
4. Güncelleme başarılı olursa version değeri 1 artırılır.
*/
func (r *PostgreWalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {

	// 1.
	query := `UPDATE wallets SET balance = $1, version = version + 1 WHERE id = $2 AND version = $3`

	// 2.
	result, err := r.getExecutor(ctx).ExecContext(
		ctx,
		query,
		wallet.Balance,
		wallet.ID,
		wallet.Version,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// 3.
	if rowsAffected == 0 {
		return domain.ErrConcurrentModification
	}

	// 4.
	// 200 ise domain nesnesini update et
	wallet.Version++

	return nil
}

/*
1. "GetIdempotencyRecord", verilen idempotency key'e ait kaydı veritabanından getirir.
2. Context'te aktif bir transaction varsa sorguyu transaction üzerinden, yoksa normal veritabanı bağlantısı üzerinden çalıştırır.
3. Gelen satırı, struct'a aktarır.
4. Kayıt bulunamadıysa hata üretmez. İlgili idempotency key'in daha önce kullanılmadığı anlamına gelir.
5. Kayıt bulunamazsa hata döndürmez..
6. nil döndürerek kaydın olmadığını belirtir.
*/
func (r *PostgreWalletRepository) GetIdempotencyRecord(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {

	// 1.
	query := `SELECT idempotency_key, response_payload, created_at FROM idempotency_records WHERE idempotency_key = $1`

	// 2.
	row := r.getExecutor(ctx).QueryRowContext(
		ctx,
		query,
		key,
	)

	var record domain.IdempotencyRecord

	// 3.
	err := row.Scan(&record.Key, &record.Response, &record.CreatedAt)
	// 4.
	if err == sql.ErrNoRows {
		return nil, nil
	}

	// 5.
	if err != nil {
		return nil, err
	}

	// 6.
	return &record, nil
}

/*
1. "SaveIdempotencyRecord", başarılı bir işlemin sonucunu idempotency tablosuna kaydeder.
2. Context'te aktif bir transaction varsa kaydı transaction üzerinden, yoksa normal veritabanı bağlantısı üzerinden ekler.
3. Aynı idempotency key ile tekrar istek geldiğinde bu kayıt kullanılarak işlemin yeniden çalıştırılması engellenir.
*/
func (r *PostgreWalletRepository) SaveIdempotencyRecord(ctx context.Context, record *domain.IdempotencyRecord) error {

	// 1.
	query := `INSERT INTO idempotency_records (idempotency_key, response_payload, created_at) VALUES ($1, $2, $3)`

	// 2.
	_, err := r.getExecutor(ctx).ExecContext(
		ctx,
		query,
		record.Key,
		record.Response,
		record.CreatedAt,
	)

	return err
}

/*
1."SaveTransaction", gerçekleştirilen para hareketini transactions tablosuna kaydeder.
2.Context'te aktif bir transaction varsa INSERT işlemini transaction üzerinden, yoksa normal veritabanı bağlantısı üzerinden gerçekleştirir.
3.İşlem sırasında oluşan hata veya başarı durumu loglanır.
*/
func (r *PostgreWalletRepository) SaveTransaction(ctx context.Context, tn *domain.Transaction) error {

	// Log amaçlı transaction ID'sini yazdırmak.
	log.Printf("SaveTransaction ID=%s", tn.ID)

	// 1
	query := `INSERT INTO transactions (id, wallet_id, amount, type, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	// 2
	_, err := r.getExecutor(ctx).ExecContext(
		ctx,
		query,
		tn.ID,
		tn.WalletID,
		tn.Amount,
		tn.Type,
		tn.Status,
		tn.CreatedAt,
	)

	// 3
	if err != nil {
		log.Printf("INSERT FAILED: %v", err)
	} else {
		log.Println("INSERT SUCCESS")
	}
	return err
}

/*
1."GetTransactionsByWalletID", verilen wallet ID'sine ait tüm işlem kayıtlarını getirir.
2.Context'te aktif bir transaction varsa sorguyu transaction üzerinden, yoksa normal veritabanı bağlantısı üzerinden çalıştırır.
3.Sonuç olarak ilgili wallet'ın tüm transaction geçmişini döndürür.
*/
func (r *PostgreWalletRepository) GetTransactionsByWalletID(ctx context.Context, walletID string) ([]*domain.Transaction, error) {
	// 1.
	query := `SELECT id, wallet_id, amount, type, status, created_at FROM transactions WHERE wallet_id = $1`

	// 2.
	rows, err := r.getExecutor(ctx).QueryContext(
		ctx,
		query,
		walletID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close() // İşimiz bitince ResultSet'i kapat.

	var transactions []*domain.Transaction

	// Dönen tüm satırları sırayla oku.
	for rows.Next() {
		var tn domain.Transaction
		// Her satırı Transaction struct'ına aktar.
		if err := rows.Scan(
			&tn.ID,
			&tn.WalletID,
			&tn.Amount,
			&tn.Type,
			&tn.Status,
			&tn.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tn) // Listeye ekle.
	}
	return transactions, nil // Tüm transaction listesini döndür.
}

// ! "BeginTx", yeni bir transaction başlatır ve transaction'ı context içine ekler.
func (r *PostgreWalletRepository) BeginTx(ctx context.Context) (context.Context, error) {

	tx, err := r.db.BeginTx(ctx, nil) // TRANSACTION başlatılır. Yani Database'e "yeni bir transaction başlat" diyoruz. "POSTGRESQL Karşılığı - BEGIN"  "tx artık *sql.Tx" tipinde bir transaction nesnesidir.
	if err != nil {
		return nil, err
	}

	// Go'daki "context.Context" immutable'dır. Mevcut "context" değişmez.
	return WithTx(ctx, tx), nil // context içine tx'i gömüyoruz
}

// ! "getExecutor", context içinde "transaction" varsa onu, yoksa normal "database" bağlantısını döndürür.
func (r *PostgreWalletRepository) getExecutor(ctx context.Context) DBExecutor {

	// Context içinden "GetTx(ctx)" ile bir transaction (tx) alınmaya çalışılıyor, Eğer tx != nil değilse yani transaction varsa "O TRANSACTION'u" döner.
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return r.db // → Eğer context'te transaction yoksa, repository'nin normal veritabanı bağlantısı (r.db) döndürülüyor.
}

// ! "Commit", transaction boyunca yapılan değişiklikleri kalıcı hale getirir.
func (r *PostgreWalletRepository) Commit(ctx context.Context) error {
	tx := GetTx(ctx) // "WithTx(ctx, tx)" ile context'e koyduğun transaction'ı geri alıyor.
	if tx == nil {
		return errors.New("no transaction found in context.") // Eğer "context" içinde "transaction" yoksa commit edilecek bir şey de yoktur. HATA DÖNER
	}

	return tx.Commit() // DB Tarafında şunu söyler. "Transaction boyunca yaptığım tüm değişiklikleri kalıcı hale getir."
}

// ! "Rollback", transaction sırasında hata oluşursa, yapılan değişiklikleri geri alır.
func (r *PostgreWalletRepository) Rollback(ctx context.Context) error {
	tx := GetTx(ctx)
	if tx == nil { // Transaction yoksa "ROLLBACK" edecek birşey yoktur.
		return nil // Hata dönmeyiz. Çünkü "rollback" güvenli şekilde temizlik (cleanup) yapar.
	}
	return tx.Rollback()
}

/*
1."UpdateTransactionStatus", verilen transaction'ın durumunu günceller.
2.Context'te aktif bir transaction varsa UPDATE işlemini transaction üzerinden, yoksa normal veritabanı bağlantısı üzerinden gerçekleştirir.
3.Güncelleme sırasında oluşan hatalar loglanır ve çağıran metoda iletilir.
*/
func (r *PostgreWalletRepository) UpdateTransactionStatus(ctx context.Context, transactionID string, status domain.TransactionStatus) error {
	// Log amaçlı güncellenecek transaction bilgisini yazdır.
	log.Printf("UPDATE STATUS id=%s status=%s", transactionID, status)
	// 1.
	query := `UPDATE transactions SET status = $1 WHERE id = $2`

	// 2.
	_, err := r.getExecutor(ctx).ExecContext(
		ctx,
		query,
		status,
		transactionID,
	)

	// 3.
	if err != nil {
		log.Printf("UPDATE FAILED: %v", err)
	}
	return err
}

/*
1.UpdateStatus, wallet'ın durumunu (status) optimistic locking kullanarak günceller.
2.Context'te aktif bir transaction varsa UPDATE işlemini transaction üzerinden, yoksa normal veritabanı bağlantısı üzerinden gerçekleştirir.
3.Version uyuşmazsa eşzamanlı güncelleme olduğu kabul edilir ve domain.ErrConcurrentModification hatası döndürülür.
*/
func (r *PostgreWalletRepository) UpdateStatus(ctx context.Context, id string, status domain.WalletStatus, currentVersion int) error {

	// 1.
	// Optimistic Locking: WHERE version = $3 ile işlem anındaki versiyonu kontrol ediyoruz. "version = version + 1:" Veritabanı seviyesinde atomik artış sağlar.
	// Başarılı olursa versiyonu veritabanında increment ediyoruz (version = version + 1).
	// WHERE version = $3: Cüzdanı kapatmaya çalışırken, birisi aynı anda o cüzdandan para çekmeye çalışırsa versiyon artacağı için bu UPDATE işlemi başarısız olur
	query := `UPDATE wallets SET status = $1, version = version + 1 WHERE id = $2 AND version = $3`

	// 2.
	result, err := r.getExecutor(ctx).ExecContext(
		ctx,
		query,
		status,
		id,
		currentVersion,
	)

	if err != nil {
		return err
	}

	// Kaç satır güncellendi, kontrol edilir.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Eğer hiçsatır güncellenmediyse, bu "wallet" başka bir işlem tarafından güncellenmiştir veya cüzdan ID'si yanlıştır. İki durumda da "OPTIMISTIC LOCK"çakışma hatası döneriz.
	if rowsAffected == 0 {
		return domain.ErrConcurrentModification
	}

	return nil
}

/*
func (r *PostgreWalletRepository) Delete(ctx context.Context, id string) error {

	query := `DELETE FROM wallets WHERE id = $1 AND balance = 0 AND status 'CLOSED'`

}*/
