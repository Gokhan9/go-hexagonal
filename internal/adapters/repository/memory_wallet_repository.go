package repository

import (
	"context"
	"errors"
	"fmt"
	"go-hexagonal/internal/core/domain"
	"sync"
)

// "MemoryWalletRepository", WalletRepository interface'ini memory üzerinden implement eder.
type MemoryWalletRepository struct {
	wallets           map[string]*domain.Wallet
	idempotencyRecord map[string]*domain.IdempotencyRecord
	transactions      map[string][]*domain.Transaction // ! YENİ EKLENDİ. Cüzdan ID'lerini anahtar olarak kullanan ve *domain.Transaction slice'larını değer olarak tutan bir transactions map'i.
	mu                sync.Mutex                       // <-- Her wallet işlemi için read/write güvenliği sağlayacak kilit!
}

// "NewMemoryWalletRepository", yeni bir memory deposu create eder.
func NewMemoryWalletRepository() *MemoryWalletRepository {
	return &MemoryWalletRepository{
		wallets:           make(map[string]*domain.Wallet),
		idempotencyRecord: make(map[string]*domain.IdempotencyRecord),
		transactions:      make(map[string][]*domain.Transaction), // ! YENİ EKLENDİ.
	}
}

// "Create" ile yeni bir wallet'i memorye kaydeder.
func (r *MemoryWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {

	r.mu.Lock()         // → Aynı anda farklı goroutine'lerin "wallets" map'ine erişimini kısıtlar. (Aynı anda sadece TEK BİR "goroutine" güncelleyebilir)
	defer r.mu.Unlock() // → function bitince kilidi açar. (RACE CONDITION önlemek.)

	// → r.wallets içerisinde "wallet.ID" ile "map" içinde arama yapıyor.
	if _, exists := r.wallets[wallet.ID]; exists {
		return errors.New("wallet already exists.")
	}

	r.wallets[wallet.ID] = wallet // ? Yeni "wallet", "Map" içine eklenir.
	return nil
}

// GetByID'ye göre cüzdanı getir
func (r *MemoryWalletRepository) GetByID(ctx context.Context, id string) (*domain.Wallet, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	wallet, exists := r.wallets[id]
	if !exists {
		return nil, errors.New("wallet not found..")
	}

	//? Go'da pointer döndüğümüz için, dışarıda yer alan katmanların(servis) map'te ki orijinal veriyi kilit dışındayken manipüle etmemesi için nesnenin kopyasını (deep copy) dönmeliyiz.
	clonedWallet := *wallet
	return &clonedWallet, nil
}

// Update ile mevcut olan wallet'i günceller..
func (r *MemoryWalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	currentWallet, exists := r.wallets[wallet.ID]
	if !exists {
		return errors.New("wallet not found")
	}

	// --- DEBUG LOGLARI 1 ---
	fmt.Printf("DEBUG: Updating wallet %s. Current Balance: %d, New Balance: %d\n",
		wallet.ID, currentWallet.Balance, wallet.Balance)

	// Servisten gelen versiyon, bende ki güncel versiyona eşit mi? değilse hatayı dön.
	if wallet.Version != currentWallet.Version {
		return domain.ErrConcurrentModification
	}

	// Doğrulama başarılıysa versiyonu 1 artırıp kaydet
	wallet.Version++

	// Map içerisine clone(kopya) saklamak güvenilir yöntemlerden biri.
	cloned := *wallet
	r.wallets[wallet.ID] = &cloned

	// --- DEBUG LOGLARI 2 ---
	fmt.Printf("DEBUG: Update successful for wallet %s. New version: %d\n",
		wallet.ID, wallet.Version)
	return nil
}

/*
- (r *MemoryWalletRepository) - Fonksiyonun MemoryWalletRepository isimli yapıya (struct) ait bir metot (receiver) olduğunu ve verileri değiştirebilmek için
işaretçi (pointer) kullandığını gösterir.
- GetIdempotencyRecord - Fonksiyona verilen, benzersiz işlem anahtarına göre arama yapacağını ifade eden isimdir.
- (ctx context.Context, key string) - Fonksiyonun girdi parametreleridir; ctx işlem sürelerini ve iptalleri yönetir, key ise aranacak benzersiz işlem anahtarı metnidir.
- (*domain.IdempotencyRecord, error) - Fonksiyonun çıktı (dönüş) değerleridir; işlem başarılıysa bulunan kaydın adresini (*domain.IdempotencyRecord), başarısızsa hata detayını
(error) döndürür.
*/
func (r *MemoryWalletRepository) GetIdempotencyRecord(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	record, exists := r.idempotencyRecord[key]
	if !exists {
		return nil, nil
	}
	return record, nil
}

func (r *MemoryWalletRepository) SaveIdempotencyRecord(ctx context.Context, record *domain.IdempotencyRecord) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.idempotencyRecord[record.Key] = record // ? Yeni "Record", "Map" içine eklenir.
	return nil
}

func (r *MemoryWalletRepository) SaveTransaction(ctx context.Context, tn *domain.Transaction) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions[tn.WalletID] = append(r.transactions[tn.WalletID], tn) // ! İşlem Kaydı eklerken "append" ile listeye ekleme yaptık.
	return nil
}

/*
Geçmiş kayıtları getirirken dışarıdan manipüle edilmemesi için slice'ın bir kopyası (deep copy) döndürülüyor.
*/
func (r *MemoryWalletRepository) GetTransactionsByWalletID(ctx context.Context, walletID string) ([]*domain.Transaction, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	/*
		tns ------> [ ptr1 ][ ptr2 ]
		cloned ---> [ ptr1 ][ ptr2 ]
	*/
	tns := r.transactions[walletID]                 // 1. walletID ile ilgili slice çektik.
	cloned := make([]*domain.Transaction, len(tns)) // Yeni Slice create edildi, NOT:SLICE KOPYALANIYOR, TRANSACTIONS OBJELERİ DEĞİL.
	copy(cloned, tns)                               // "copy" ile elemanlar içine aktarılıyor.
	return cloned, nil
}
