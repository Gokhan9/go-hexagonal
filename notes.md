internal/adapters/handler/wallet_handler.go

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