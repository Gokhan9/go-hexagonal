# Go Hexagonal Wallet API

Bu proje, Hexagonal Mimari (Ports & Adapters) prensiplerine uygun olarak geliştirilmiş bir Finansal Cüzdan API'sidir.

## Özellikler
- **Hexagonal Mimari:** Domain, Ports ve Adapters ayrımı.
- **Finansal Güvenlik:** ACID uyumlu transaction yönetimi.
- **Audit Loglama:** Tüm kritik finansal işlemler (Transfer, Deposit vb.) `audit_logs` tablosunda kayıt altına alınır.
- **İşlem Durum Yönetimi:** `PENDING`, `COMPLETED`, `FAILED` statüleri ile işlem takibi.
- **İdempotency:** `X-Idempotency-Key` desteği ile mükerrer istek koruması.
- **Optimistic Locking:** Eşzamanlı bakiye güncellemelerinde veri tutarlılığı.
- **Veritabanı:** PostgreSQL.
- **Docker:** `docker-compose.yml` ile hızlı kurulum.

## Kurulum ve Çalıştırma

### Gereksinimler
- Docker & Docker Compose
- Go 1.23+

### Başlatma
Proje kök dizininde aşağıdaki komutla veritabanını ve uygulamayı ayağa kaldırabilirsiniz:

```bash
docker-compose up -d
```

### Test
Smoke testlerini çalıştırmak için `backend` dizinine gidip şu komutu kullanabilirsiniz:
```bash
cd backend
go test -v ./internal/test/smoke_test.go
```

## API Kullanımı
- **Cüzdan Oluştur:** `POST /wallets`
- **Para Yatır:** `POST /wallets/{id}/deposit`
- **Para Çek:** `POST /wallets/{id}/withdraw`
- **Transfer:** `POST /wallets/{id}/transfer`
- **Bakiye Sorgula:** `GET /wallets/{id}/balance`
