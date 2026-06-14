📌 PORTS:

1 - Driver/Primary(Birincil) PORT ve Driven/Secondary(İkincil) PORT olmak üzere 2 adet Port'u kullanıyoruz. 
- Driver/Primary(Birincil) PORT: Handler(service) ile ilgili işlemler yapmak istiyorsak sahip olduğumuz Handler'lar bu interface'i implement edecekler.(WalletService)
- Driven/Secondary(İkincil) PORT: Application'un veriyi nasıl saklayacağını bilir. Ayrıca DB(postgres,redis vb..) bu interface'i implement eder.(WalletRepository)

📌 ADAPTERS: 

2 - Portları harici bileşenlerle bağlayan parça adaptördür. İki tür adaptör vardır.
- Driver Adapter / Primary Adapter: Business logic işlemini gerçekleştirmek için "PORT" interface'ini kullanır..
- Driven Adapter / Secondary Adapter: Uygulama, dış bileşenlerle(databaseler(postgre,mongo), dış servisler(Ödeme servisleri(iyzico)), API'ler(hava durumu API'si), Mesajlaşma sistemi(rabbitmq, kafka), e-posta,bildirim sistemleri(SMTP, Firebase Cloud Message), Harici SaaS Sistemleri(CRM,ERP)) iletişim kurmak için bunu kullanır. İş mantığının isteğini, dış teknoloji bileşenlerinin isteklerine dönüştürür.








🔥 Örnek

Domain error: var ErrorInsufficientFunds = errors.New("insufficient funds")

🛡️ TEST 

→ TestWalletService_Deposit_InvalidAmount
→ TestWalletService_Withdraw_In_SufficientFunds

1. Guard Clause (Koruyucu Koşul) Tasarımı: Metodun asıl ağır iş yüküne (veritabanından cüzdanı çekmek, kilitlemek vb.) girmeden önce girdilerin doğruluğunu en başta 
kontrol etmek (fail-fast) performansı artırır ve gereksiz DB/Memory yükünü engeller.
2. Domain Güvenliği: Finansal dünyada sıfır veya negatif bakiye işlemleri dolandırıcılığa (exploit) en açık yerlerdir. Negatif bir değer gönderildiğinde balance = balance + (-100) işlemi çalışarak Deposit fonksiyonunun gizlice bir Withdraw işlemine dönüşmesini engellemiş olduk.
3. Mimaride Sorumluluk Dağılımı (Separation of Concerns): Validasyon hata tanımları Domain katmanında bulunur, çünkü bu hata iş mantığının bir parçasıdır. Service katmanı ise bu kuralı uygular.


→ TestWalletService_Concurrent_Deposit 

-wallet.go içerisine "Version int" eklendi.
-memory_wallet_repository.go içerisinde yer alan MemoryWalletRepository struct değeri "sync.Mutex" çevrildi. Ayrıca "Update" fonksiyonunda güncelleme yapıldı.
-wallet_service.go içerisinde ki Deposit fonksiyonunda ki kod bloğu for döngüsüne alındı.

- İyimser Kilitleme (Optimistic Locking) — DDD 
→ Gerçek projelerde (PostgreSQL/MySQL kullanırken) servis katmanına Mutex koymak performansı düşürür ve birden fazla sunucu (Replica/Pod) çalıştığında işe yaramaz. Bunun yerine nesneye bir "Version" alanı eklenir. Veritabanına güncellenmiş nesne gönderilirken "Eğer bendeki versiyon hâlâ veritabanındakiyle aynıysa güncelle" denir. Eğer başkası araya girip versiyonu değiştirdiyse hata fırlatılır ve işlem yeniden denenir (Retry).





✏️ test/wallet_service_test.go

→ repo := repository.NewMemoryWalletRepository() → Bellek içi(In-Memory) repo oluşturmak. DB bağımlılığı yok. "Postgre,Mongo" gibi db araçlarını kullanmayız.. Yerine "MemoryWalletRepository".
- Veriler ram'de tutulur.

→ service := services.NewWalletService(repo) → Service
→ ctx := context.Background() → boş context.

→ require.NoError(t, err) → Hata olması durumunda testi durdurur.

"require.NoError", aşağıda ki yapıya benzer bir hata döner.
if err != nil {
    t.FailNow()
}

→ require.ErrorIs(t, err, domain.ErrorInsufficientFunds) → "Error mu yoksa Wrap edilmiş mi onu kontrol eder"
- Error'un "WRAP" edilmiş olup olmaması, hatanın başka bir yapı(wrapper) içinde sarılıp/sarılmadığını anlatır.
Not: "require.ErrorIs", error zincirini gezer ve wrapped errorları kontrol eder. İçeride tanımladığımız "ERROR" var mı yok mu onu kontrol eder.

- var wg sync.WaitGroup → Add, done ve wait işlemlerini başlatmak..
→ wg.Add(goroutineCount) → Beklenecek goroutine sayısı(örn:5)
→ defer wg.Done() → goroutine işini tamamladığında sayaçtan "1" eksilir. (Function içinde en başa yazılır.)
→ wg.Wait() → Sayaç 0 olana kadar diğer işlemleri bloklarız, 0 olduğunda program kaldığı yerden devam edebilir.



✏️ handler/wallet_handler.go

→ "http.ResponseWriter", Go’da HTTP response (sunucu cevabı) yazmak için kullanılan bir arayüzdür (interface).
w → response (cevap yazacağın yer)
r → request (istek bilgisi)


→ json.NewDecoder(r.body).Decode(&req) → Client'tan gelen HTTP Request'in body'sine bakar.(API üzerinden gönderilen) "JSON" formatında ki veriyi
doğrudan Go içerisinde ki bir "struct'a dönüştürür(parse/decode)"
→ h.service.CreateWallet(r.Context(), req.Owner, req.Currency) → Service'a bağlı "CreateWallet" fonksiyonunu, istek bağlamı(context, cüzdan sahibi(Owner) ve
para birimi(Currency)) parametreleriyle çağır. Dönen sonuçları "wallet ve err" değişkenlerine ata.
→ id := r.PathValue("id") → "Path" parametre okuma.
→ json":"created_at" not compatible with reflect.StructTag.Get: bad syntax for struct tag pair


✏️ repository/memory_wallet_repository.go

→ wallets map[string]*domain.Wallet → "String anahtar" → Wallet pointer değeri tutan map
→ sync.Mutex → Her wallet işlemi için read/write güvenliği sağlayacak kilit!
→ r.mu.Lock() → Yazma(Write) Kilidi (Aynı anda sadece TEK BİR "goroutine" güncelleyebilir, Yazma(Write) işlemine başlamadan önce kilitle..)
→ defer r.mu.Unlock() → function bitince kilidi açar. (RACE CONDITION önlemek.)









