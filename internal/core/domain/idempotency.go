package domain

import "time"

// ! YENİ EKLENDİ.
type IdempotencyRecord struct {
	Key       string    // Client'tan gelen Unique Key
	Response  []byte    // İşlem sonucu return edilen response'un JSON hali
	CreatedAt time.Time // Kayıt Zamanı
}
