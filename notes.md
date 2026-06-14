PORTS:

1 - Driver/Primary(Birincil) PORT ve Driven/Secondary(İkincil) PORT olmak üzere 2 adet Port'u kullanıyoruz. 
- Driver/Primary(Birincil) PORT: Handler(service) ile ilgili işlemler yapmak istiyorsak sahip olduğumuz Handler'lar bu interface'i implement edecekler.(WalletService)
- Driven/Secondary(İkincil) PORT: Application'un veriyi nasıl saklayacağını bilir. Ayrıca DB(postgres,redis vb..) bu interface'i implement eder.(WalletRepository)

ADAPTERS: 

2 - Portları harici bileşenlerle bağlayan parça adaptördür. İki tür adaptör vardır.
- Driver Adapter / Primary Adapter: Business logic işlemini gerçekleştirmek için "PORT" interface'ini kullanır..
- Driven Adapter / Secondary Adapter: Uygulama, dış bileşenlerle(databaseler(postgre,mongo), dış servisler(Ödeme servisleri(iyzico)), API'ler(hava durumu API'si), Mesajlaşma sistemi(rabbitmq, kafka), e-posta,bildirim sistemleri(SMTP, Firebase Cloud Message), Harici SaaS Sistemleri(CRM,ERP)) iletişim kurmak için bunu kullanır. İş mantığının isteğini, dış teknoloji bileşenlerinin isteklerine dönüştürür.




internal/adapters/handler/wallet_handler.go

→ "http.ResponseWriter", Go’da HTTP response (sunucu cevabı) yazmak için kullanılan bir arayüzdür (interface).
w → response (cevap yazacağın yer)
r → request (istek bilgisi)


→ json.NewDecoder(r.body).Decode(&req) → Client'tan gelen HTTP Request'in body'sine bakar.(API üzerinden gönderilen) "JSON" formatında ki veriyi
doğrudan Go içerisinde ki bir "struct'a dönüştürür(parse/decode)"

→ h.service.CreateWallet(r.Context(), req.Owner, req.Currency) → Service'a bağlı "CreateWallet" fonksiyonunu, istek bağlamı(context, cüzdan sahibi(Owner) ve
para birimi(Currency)) parametreleriyle çağır. Dönen sonuçları "wallet ve err" değişkenlerine ata.

→ id := r.PathValue("id") → "Path" parametre okuma.

→ json":"created_at" not compatible with reflect.StructTag.Get: bad syntax for struct tag pair



test/wallet_service_test.go

→ repo := repository.NewMemoryWalletRepository() → Bellek içi(In-Memory) repo oluşturmak. DB bağımlılığı yok. "Postgre,Mongo" gibi db araçlarını kullanmayız.. Yerine "MemoryWalletRepository".
- Veriler ram'de tutulur.

→ service := services.NewWalletService(repo) → Service
→ ctx := context.Background() → boş context.

→ require.NoError(t, err) → Hata olması durumunda testi durdurur.

if err != nil {
    t.FailNow()
}

→ require.ErrorIs(t, err, domain.ErrorInsufficientFunds) → "Error mu yoksa Wrap edilmiş mi onu kontrol eder"
- Error'un "WRAP" edilmiş olup olmaması, hatanın başka bir yapı(wrapper) içinde sarılıp/sarılmadığını anlatır.
Not: "require.ErrorIs", error zincirini gezer ve wrapped errorları kontrol eder. İçeride tanımladığımız "ERROR" var mı yok mu onu kontrol eder.

🔥 Örnek

Domain error: var ErrorInsufficientFunds = errors.New("insufficient funds")

🛡️ TEST : TestWalletService_Deposit_InvalidAmount ve TestWalletService_Withdraw_In_SufficientFunds

1. Guard Clause (Koruyucu Koşul) Tasarımı: Metodun asıl ağır iş yüküne (veritabanından cüzdanı çekmek, kilitlemek vb.) girmeden önce girdilerin doğruluğunu en başta 
kontrol etmek (fail-fast) performansı artırır ve gereksiz DB/Memory yükünü engeller.
2. Domain Güvenliği: Finansal dünyada sıfır veya negatif bakiye işlemleri dolandırıcılığa (exploit) en açık yerlerdir. Negatif bir değer gönderildiğinde balance = balance + (-100) işlemi çalışarak Deposit fonksiyonunun gizlice bir Withdraw işlemine dönüşmesini engellemiş olduk.
3. Mimaride Sorumluluk Dağılımı (Separation of Concerns): Validasyon hata tanımları Domain katmanında bulunur, çünkü bu hata iş mantığının bir parçasıdır. Service katmanı ise bu kuralı uygular.




