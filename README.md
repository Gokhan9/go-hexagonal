# Go Hexagonal Wallet Service

![Domain Driven Hexagon Architecture](assets/DomainDrivenHexagon.png)

Bu proje, Go dilinde **Hexagonal Architecture (Ports and Adapters)** prensipleri kullanılarak geliştirilmiş modern bir cüzdan (wallet) servisidir.

## 🏗 Mimari Yapı

Proje, bağımlılıkların içe doğru (core'a doğru) olduğu, iş mantığının dış dünyadan (DB, API, CLI) tamamen izole edildiği bir yapıdadır.

- **Internal/Core/Domain:** Uygulamanın kalbi. Cüzdan ve İşlem (Transaction) modelleri ile temel iş kuralları (Balance kontrolü vb.) burada yer alır.
- **Internal/Core/Ports:** Uygulamanın dış dünya ile iletişim kontratları (Interface'ler). 
    - *Driving Ports:* Servis yetenekleri.
    - *Driven Ports:* Veri saklama (repository) ihtiyaçları.
- **Internal/Core/Services:** İş mantığının (Use-Case) koordine edildiği katman.
- **Internal/Adapters:**
    - **Inbound (Driving):** HTTP Handlers (REST API).
    - **Outbound (Driven):** Memory Repository (Persistence).

## 🚀 Başlangıç

### Gereksinimler
- Go 1.25.4 veya üzeri

### Kurulum ve Çalıştırma
1. Bağımlılıkları yükleyin:
   ```bash
   go mod download
   ```
2. Uygulamayı ayağa kaldırın:
   ```bash
   go run cmd/main.go
   ```

## 🛠 Özellikler
- [x] Cüzdan Oluşturma
- [x] Cüzdan Bakiyesi Sorgulama
- [x] Para Yatırma (Deposit)
- [x] Para Çekme (Withdraw)
- [ ] PostgreSQL Entegrasyonu (Gelecek Plan)
- [ ] Unit & Integration Testleri (Gelecek Plan)

## 📈 Proje Durumu
Şu an **Phase 4: Infrastructure Adapters** aşamasındayız. Temel domain yapısı, portlar ve bellek tabanlı repository tamamlandı.
