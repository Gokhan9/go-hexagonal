PORTS:

1 - Driver/Primary(Birincil) PORT ve Driven/Secondary(İkincil) PORT olmak üzere 2 adet Port'u kullanıyoruz. 
- Driver/Primary(Birincil) PORT: Handler(service) ile ilgili işlemler yapmak istiyorsak sahip olduğumuz Handler'lar bu interface'i implement edecekler.(WalletService)
- Driven/Secondary(İkincil) PORT: Application'un veriyi nasıl saklayacağını bilir. Ayrıca DB(postgres,redis vb..) bu interface'i implement eder.(WalletRepository)

ADAPTERS: 

2 - Portları harici bileşenlerle bağlayan parça adaptördür. İki tür adaptör vardır.
- Driver Adapter / Primary Adapter: Business logic işlemini gerçekleştirmek için "PORT" interface'ini kullanır..
- Driven Adapter / Secondary Adapter: Uygulama, dış bileşenlerle(databaseler(postgre,mongo), dış servisler(Ödeme servisleri(iyzico)), API'ler(hava durumu API'si), Mesajlaşma sistemi(rabbitmq, kafka), e-posta,bildirim sistemleri(SMTP, Firebase Cloud Message), Harici SaaS Sistemleri(CRM,ERP)) iletişim kurmak için bunu kullanır. İş mantığının isteğini, dış teknoloji bileşenlerinin isteklerine dönüştürür.







nternal/adapters/handler/wallet_handler.go

/*

→ "http.ResponseWriter", Go’da HTTP response (sunucu cevabı) yazmak için kullanılan bir arayüzdür (interface).
w → response (cevap yazacağın yer)
r → request (istek bilgisi)

*/

/*

json.NewDecoder(r.body).Decode(&req) → Client'tan gelen HTTP Request'in body'sine bakar.(API üzerinden gönderilen) "JSON" formatında ki veriyi
doğrudan Go içerisinde ki bir "struct'a dönüştürür(parse/decode)"

h.service.CreateWallet(r.Context(), req.Owner, req.Currency) → Service'a bağlı "CreateWallet" fonksiyonunu, istek bağlamı(context, cüzdan sahibi(Owner) ve
para birimi(Currency)) parametreleriyle çağır. Dönen sonuçları "wallet ve err" değişkenlerine ata.

*/

/*
→ id := r.PathValue("id") // * "Path" parametre okuma.
*/


/*
json":"created_at" not compatible with reflect.StructTag.Get: bad syntax for struct tag pair
*/



