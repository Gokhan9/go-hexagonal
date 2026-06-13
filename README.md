# Go Hexagonal Wallet Service

![Domain Driven Hexagon Architecture](DomainDrivenHexagon.png)

Bu proje, Go dilinde **Hexagonal Architecture (Ports and Adapters)** prensipleri kullanılarak geliştirilmiş modern bir cüzdan (wallet) servisidir.

## 🏗 Mimari Yapı

Proje, bağımlılıkların içe doğru (core'a doğru) olduğu, iş mantığının dış dünyadan (DB, API, CLI) tamamen izole edildiği bir yapıdadır.

- **Internal/Core/Domain:** Uygulamanın kalbi. Cüzdan ve İşlem (Transaction) modelleri ile temel iş kuralları (Balance kontrolü, yetersiz bakiye doğrulaması) burada yer alır.
- **Internal/Core/Ports:** Uygulamanın dış dünya ile iletişim kontratları (Arayüzler / Interface'ler). 
    - *Driving Ports (Birincil Liman):* HTTP API gibi tetikleyicilerin çağırdığı servis kontratları (`WalletService`).
    - *Driven Ports (İkincil Liman):* Uygulamanın veriyi saklamak için ihtiyaç duyduğu veri tabanı/bellek kontratları (`WalletRepository`).
- **Internal/Core/Service:** İş mantığının (Use-Case) koordine edildiği ve portların uygulandığı katman.
- **Internal/Adapters:**
    - **Inbound (Driving - Girdi Adaptörleri):** Go 1.22+ yerleşik HTTP yönlendiricisi (`http.NewServeMux`) kullanılarak yazılmış HTTP Handlers (`internal/adapters/handler`).
    - **Outbound (Driven - Çıktı Adaptörleri):** Thread-safe bellek tabanlı veri deposu (`internal/adapters/repository`).
- **Internal/Api/Dto:** API istek ve yanıt şablonları (`internal/api/dto`). Kuruş (int64) ve ana birim (float64) dönüşümlerini bu katman yönetir.

## 🚀 Başlangıç

### Gereksinimler
- Go 1.25 veya üzeri

### Kurulum ve Çalıştırma
1. Projeyi klonlayın ve bağımlılıkları yükleyin:
   ```bash
   go mod download
   ```
2. Uygulamayı ayağa kaldırın:
   ```bash
   go run cmd/main.go
   ```

## 🛠 API Uç Noktaları (Endpoints)

Go 1.22'nin yeni yönlendirme yetenekleri kullanılarak tanımlanmış endpoint listesi ve örnek `curl` istekleri:

### 1. Yeni Cüzdan Oluşturma (POST `/wallets`)
```bash
curl -X POST http://localhost:8080/wallets \
  -H "Content-Type: application/json" \
  -d '{"owner": "Gökhan", "currency": "TRY"}'
```

### 2. Cüzdan Bilgilerini Getirme (GET `/wallets/{id}`)
```bash
curl -X GET http://localhost:8080/wallets/<wallet_id>
```

### 3. Para Yatırma (POST `/wallets/{id}/deposit`)
```bash
curl -X POST http://localhost:8080/wallets/<wallet_id>/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount": 150.75}'
```

### 4. Para Çekme (POST `/wallets/{id}/withdraw`)
```bash
curl -X POST http://localhost:8080/wallets/<wallet_id>/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount": 50.25}'
```

## 🛠 Özellikler
- [x] Cüzdan Oluşturma
- [x] Cüzdan Bakiyesi Sorgulama
- [x] Para Yatırma (Deposit - Float to Cents otomatik dönüşümüyle)
- [x] Para Çekme (Withdraw - Yetersiz bakiye kontrolleriyle)
- [ ] PostgreSQL Entegrasyonu (Gelecek Plan)
- [ ] Unit & Integration Testleri (Gelecek Plan)

## 📈 Proje Durumu
Şu an **Phase 5: Input Adapters & API Endpoints** başarıyla tamamlandı. API uç noktaları, Go'nun yerleşik `net/http` paketi kullanılarak, herhangi bir üçüncü parti framework (Gin, Chi vb.) bağımlılığı olmaksızın enjekte edildi ve başarıyla çalışmaktadır.
