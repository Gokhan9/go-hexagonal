# 🚀 Go Hexagonal Wallet Service - Geliştirici Notları

Bu doküman, projede uygulanan **Hexagonal Mimari (Ports & Adapters)**, **Eşzamanlılık (Concurrency)** yönetimi, **Mükerrer İşlem Koruması (Idempotency)** ve test stratejilerine dair teknik detayları ve tasarım kararlarını içerir.

---

## 📌 İçindekiler
1. [Mimaride Portlar ve Adaptörler (Ports & Adapters)](#1-mimaride-portlar-ve-adaptörler-ports--adapters)
2. [Güvenli ve Sağlam Domain Tasarımı](#2-güvenli-ve-sağlam-domain-tasarımı)
3. [Eşzamanlılık Kontrolü ve Yarış Durumları (Concurrency & Race Conditions)](#3-eşzamanlılık-kontrolü-ve-yarış-durumları-concurrency--race-conditions)
4. [Idempotency (Mükerrer İşlem Koruması)](#4-idempotency-mükerrer-işlem-koruması)
5. [Katmanların Detaylı İncelenmesi](#5-katmanların-detaylı-incelenmesi)
   - [Test Katmanı (`wallet_service_test.go`)](#test-katmanı-wallet_service_testgo)
   - [Sunum Katmanı (`wallet_handler.go`)](#sunum-katmanı-wallet_handlergo)
   - [Veri Erişim Katmanı (`memory_wallet_repository.go`)](#veri-erişim-katmanı-memory_wallet_repositorygo)

---

## 1. Mimaride Portlar ve Adaptörler (Ports & Adapters)

### 🔌 Portlar (Interfaces)
Uygulamada iki tür port kullanılır:

*   **Driver / Primary (Birincil) PORT:** Sunum katmanının (örn. HTTP API) iş mantığına erişebilmesi için tanımlanan arayüzdür. İş yeteneklerinin kontratıdır.
    *   *Örnek:* `WalletService` arayüzü. API Handler'ları bu portu çağırır.
*   **Driven / Secondary (İkincil) PORT:** İş mantığının dış dünyaya (veri saklama, dış servisler) veri göndermek/çekmek için kullandığı kontratlardır.
    *   *Örnek:* `WalletRepository` arayüzü. Veritabanı (PostgreSQL, Redis vb.) bu arayüzü implement eder.

### ⚙️ Adaptörler (Adapters)
Portları harici bileşenlere bağlayan somut uygulama sınıflarıdır:

*   **Driver Adapter / Primary Adapter:** İş mantığı işlemlerini tetiklemek için birincil portu kullanır.
    *   *Örnek:* HTTP API Handler katmanı. Gelen HTTP isteklerini ayrıştırıp `WalletService` portunu tetikler.
*   **Driven Adapter / Secondary Adapter:** İş mantığının taleplerini dış teknolojilerin diline dönüştürür ve ikincil portları uygular.
    *   *Örnek:* Veritabanları (PostgreSQL, MongoDB), dış ödeme entegrasyonları (İyzico), mesaj kuyrukları (RabbitMQ, Kafka), bildirim sistemleri (SMTP, Firebase) veya in-memory saklayıcılar.

---

## 2. Güvenli ve Sağlam Domain Tasarımı

Projede uygulanan domain kurallarında şu üç temel yaklaşım benimsenmiştir:

1.  **Guard Clause (Koruyucu Koşul) Tasarımı:** Metodun asıl ağır iş yüküne (veritabanından veri çekmek, kilitlemek vb.) girmeden önce girdilerin doğruluğunu en başta kontrol etmesidir (*fail-fast*). Bu yaklaşım performansı artırır ve gereksiz DB/Memory yükünü engeller.
2.  **Domain Güvenliği:** Finansal işlemlerde negatif bakiye işlemleri dolandırıcılığa (*exploit*) en açık noktalardır. Negatif bir değer gönderildiğinde `balance = balance + (-amount)` şeklinde çalışıp `Deposit` işleminin gizlice bir para çekmeye (`Withdraw`) dönüşmesi engellenmiştir.
3.  **Mimaride Sorumluluk Dağılımı (Separation of Concerns):** Validasyon hata tanımları (`ErrorInsufficientFunds`, `ErrorInvalidAmount`) Domain katmanında (`domain/errors.go`) yer alır çünkü bu hatalar doğrudan iş mantığı kuralıdır. Servis katmanı ise sadece bu kuralları uygulayıp akışı yönetir.

---

## 3. Eşzamanlılık Kontrolü ve Yarış Durumları (Concurrency & Race Conditions)

### 🚨 Yarış Durumu (Race Condition) Senaryosu
Kullanıcının bakiyesi **1000 TRY** olsun. Milisaniyeler farkla iki farklı istek gelsin:
1.  **1. İstek:** 500 TRY para çekme (Withdraw)
2.  **2. İstek:** 300 TRY para çekme (Withdraw)

Eğer bu istekler eşzamanlı çalışır ve bakiye kontrolünü aynı anda yaparlarsa, ikisi de bakiyeyi **1000 TRY** olarak görür. İki istek de onaylanır:
*   1. istek bakiyeyi günceller: **500 TRY** yapar.
*   2. istek bakiyeyi günceller (ve üstüne yazar): **700 TRY** yapar.
*   **Olması Gereken:** $1000 - 500 - 300 = 200\text{ TRY}$ bakiye kalmalıydı. Ancak bakiye 700 TRY kaldı ve sistem açık verdi (*Lost Update* problemi).

### 🛡️ Çözüm Yöntemleri

#### A. Bellek Düzeyinde Kilitleme (Pessimistic Locking / Mutex)
Bellek içi işlemlerde eşzamanlılığı korumak amacıyla `MemoryWalletRepository` üzerinde `sync.Mutex` kullanılarak yazma kilitleri (`Lock` / `Unlock`) eklenmiştir. Bu sayede aynı anda yalnızca tek bir goroutine güncellenmiş map değerine yazma yapabilir.

#### B. İyimser Kilitleme (Optimistic Concurrency Control - OCC)
Gerçek dünyada (PostgreSQL/MySQL gibi ilişkisel veritabanlarında) servis katmanına Mutex koymak performansı düşürür ve uygulamanın birden fazla sunucu (Replica/Pod) üzerinde çalıştığı dağıtık mimarilerde işe yaramaz. 

Bunun yerine **DDD (Domain-Driven Design)** prensiplerine uygun olarak nesnelere bir `Version` alanı eklenir:
1.  Nesne veritabanından mevcut versiyonuyla çekilir (`Version = 1`).
2.  İşlemler uygulanır, bakiye güncellenir.
3.  Veritabanına güncellenmiş nesne gönderilirken:
    ```sql
    UPDATE wallets SET balance = 500, version = 2 WHERE id = 'xxx' AND version = 1;
    ```
4.  Eğer başka bir işlem araya girip versiyonu `2` yaptıysa, etkilenen satır sayısı `0` döner. Bu durumda sistem hata fırlatır, servis katmanı bu eşzamanlılık hatasını yakalar, güncel nesneyi (ve yeni versiyonunu) tekrar çekip işlemi **yeniden dener (Retry)**.

> **Uygulama Notu:** `TestWalletService_Concurrent_Deposit` test senaryosunda, `wallet_service.go` içerisindeki `Deposit` fonksiyonu bu kurala uygun olarak sonsuz bir `for` döngüsüne alınmış ve başarılı güncelleme yapılana dek retry mekanizması işletilmiştir.

---

## 4. Idempotency (Mükerrer İşlem Koruması)

Ağ kopukluğu veya istemcinin butona yanlışlıkla üst üste basması nedeniyle aynı finansal isteğin sisteme mükerrer olarak gelmesi durumunda, işlemin **yalnızca bir kez** gerçekleşmesini sağlama özelliğidir.

### 🛠️ Idempotency Akış Adımları

1.  **İstemci Rolü:** İstemci (Client) her benzersiz işlem isteği için bir `Idempotency-Key` (genellikle UUID) oluşturur ve bunu HTTP isteğinin Header bilgisinde (`X-Idempotency-Key`) gönderir.
2.  **Veritabanı Kaydı (`idempotency.go`):** Gelen anahtar ve bu işleme ait işlem sonucu (`Response` verisi) eşleştirilerek saklanır.
3.  **Repository Arayüzü:** `WalletRepository` arayüzüne `GetIdempotencyRecord` ve `SaveIdempotencyRecord` kontratları eklenmiştir.
4.  **In-Memory Implementasyonu:** `MemoryWalletRepository` bileşenine `idempotencyRecords map[string]*domain.IdempotencyRecord` alanı eklenerek bellek üzerinde kayıtlar tutulmuş ve `Mutex` ile thread-safe hale getirilmiştir.
5.  **Service Katmanı Kontrolü:**
    *   Eğer gelen `idempotencyKey` boş değilse, repository'den bu key sorgulanır.
    *   Eğer bu anahtarla daha önce başarılı bir işlem yapılmışsa (**Duplicate Request**), işlem yeniden çalıştırılmaz; doğrudan önceki başarılı sonuç `nil` dönülerek (istemciye işlemin başarılı olduğu bildirilerek) fonksiyon sonlandırılır.
    *   Eşzamanlılık ve iş kuralları başarıyla tamamlandıktan sonra, yeni anahtar ve işlem sonucu veritabanına kaydedilir.

---

### Defter-i Kebir / İşlem Geçmişi (Ledger / Transaction History)

Repository içerisinde ki "GetTransactionsByWalletID" fonksiyonunun görsel açıklaması..

Diyelim repository içinde şu veri var:
r.transactions["wallet-1"] = []*domain.Transaction{
	{ID: "tx1", Amount: 100},
	{ID: "tx2", Amount: 200},
}

Fonksiyon çalışınca:
tns := r.transactions["wallet-1"]

BELLEKTE:

       tns
        ↓
+-------+-------+
| ptr A | ptr B |
+-------+-------+

ptr A → {ID:"tx1", Amount:100}
ptr B → {ID:"tx2", Amount:200}


SONRASI:
cloned := make([]*domain.Transaction, len(tns))
copy(cloned, tns)

Yeni bir slice oluşuyor:

       tns                  cloned
        ↓                      ↓
+-------+-------+      +-------+-------+
| ptr A | ptr B |      | ptr A | ptr B |
+-------+-------+      +-------+-------+
      ↘                   ↙
       {ID:"tx1", Amount:100}

      ↘                   ↙
       {ID:"tx2", Amount:200}

Yani:

tns → farklı slice
cloned → farklı slice
ama içindeki pointer’lar aynı objeleri gösteriyor.



## 5. Katmanların Detaylı İncelenmesi

### 🧪 Test Katmanı (`test/wallet_service_test.go`)

Veritabanı bağımlılığı olmadan testleri izole ve hızlı bir şekilde koşturmak için in-memory repository yapısı kullanılır:

*   **In-Memory Repo:** `repository.NewMemoryWalletRepository()` ile başlatılan sahte repository, verileri RAM'de tutar.
*   **Hata Yönetimi ve Kontrolü:** 
    *   `require.NoError(t, err)`: Beklenmeyen bir hata oluşursa testi anında durdurur (`t.FailNow()`).
    *   `require.ErrorIs(t, err, domain.ErrorInsufficientFunds)`: Hataların sarmalanmış (*wrapped*) olup olmadığını kontrol etmek için hata zincirini gezer ve asıl hatanın doğru türde olup olmadığını sorgular.
*   **Eşzamanlılık Testleri (`sync.WaitGroup`):**
    *   `var wg sync.WaitGroup`: Eşzamanlı goroutine'leri koordine etmek için kullanılır.
    *   `wg.Add(goroutineCount)`: Beklenecek toplam goroutine sayısını ayarlar.
    *   `defer wg.Done()`: Goroutine işini tamamladığında sayacı `1` azaltır.
    *   `wg.Wait()`: Tüm goroutine'ler tamamlanana (sayaç sıfırlanana) kadar ana akışı durdurur.

### 🌐 Sunum Katmanı (`handler/wallet_handler.go`)

HTTP isteklerinin karşılandığı ve yanıtlandığı katmandır:

*   `http.ResponseWriter`: Go'da istemciye yanıt yazmak (gövde, durum kodu vb.) için kullanılan standart arayüzdür.
*   `json.NewDecoder(r.Body).Decode(&req)`: İstemciden JSON formatında gelen ham HTTP gövdesini (Request Body) okuyarak Go struct yapısına dönüştürür.
*   `r.PathValue("id")`: Go 1.22+ ile gelen native router özelliği sayesinde URL yolundaki parametreleri (örneğin cüzdan ID'sini) okur.
*   `r.Header.Get("X-Idempotency-Key")`: İstemciden mükerrer işlem koruması için gönderilen benzersiz anahtarı okur ve servis katmanına iletir.

### 💾 Veri Erişim Katmanı (`memory_wallet_repository.go`)

Bellek içi hızlı veri saklama katmanıdır:

*   `wallets map[string]*domain.Wallet`: Benzersiz string anahtarlar üzerinden cüzdan pointer nesnelerini tutan eşleme (map) yapısıdır.
*   `sync.Mutex`: Eşzamanlı okuma/yazma güvenliğini (thread-safety) sağlayan kilit mekanizmasıdır.
*   `r.mu.Lock()` / `defer r.mu.Unlock()`: Yazma işlemine başlamadan önce kilit atılır, fonksiyon sonlandığında (hata olsa dahi `defer` sayesinde) kilit otomatik olarak serbest bırakılır. Bu sayede yarış durumları engellenir.
