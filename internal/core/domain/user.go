package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

// AuthMiddleware
type AuthenticatedUser struct {
	UserID   string
	Username string
}

/*
-HASH HASH-
* User'ın hesabını güvene almak, şifre sıfırlamak, yeni oluşturulan hesaplara geçici/ilk şifre atamak için "setPassword" yazarız.

- Şifre Sıfırlama : User "Şifremi Unuttum" dediğinde ve e-postasına gönderilen tek kullanımlık link/kod ile yeni bir şifre belirlemek istediğinde.
- İlk Hesap Kurulumu : Yönetici tarafından sisteme manuel eklenen veya üçüncü parti bir provider(sağlayıcı) (Google vb.) yerine e-posta/şifre ile üye olan kullanıcının
ilk şifresini kaydetmek.
- Güvenlik / Şifre Değişikliği : Kullanıcının panel üzerinden mevcut şifresini daha güncel ve güçlü bir şifre ile değiştirmek istediğinde.

! - SetPassword() → şifre oluşturma/değiştirme işlemi yapıyor ("hash" üretip kaydediyor)

NOT: Şifreler hiçbir zaman düz metin(plaintext) olarak db'ye kaydedilmemeli. "setPassword" fonksiyonunun arka planda bu parolayı/şifreyi güvenli bir şekilde hash algoritmalarıyla
(örn:bcrypt) şifreleyerek db'ye aktardığına emin olmalıyız.
*/
func (u *User) SetPassword(password string) error {

	/*
		- Girilen şifreyi "bcrypt" algoritması ile hashliyor.
		- "[]byte(password)" → parametre içinde ki "string" değeri, "bytes" dizisine çeviriyor.
		- "bcrypt.DefaultCost" → hash işleminin maliyet (güvenlik / işlem süresi) seviyesini kullanıyor. 10
		- Sonuç : bytes → oluşturulan hash,    err → hatayı döner.
		- "u.PasswordHash = string(bytes)" oluşturulan hash'i "USER" nesnesinin "PasswordHash" alanına kaydediyor.
	*/
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(bytes)
	return nil
}

/*
* User'ın girdiği şifrenin sistemde kayıtlı olan şifreyle güvenli bir şekilde eşleşip, eşleşmediğini doğrulamak için yazılır.

- Güvenlik(Hashing) : Şifreler, db'de (plaintext) saklanamaz. "CheckPassword", db'den gelen hashlenmiş şifre ile user'ın giriş yaparken yazdığı şifreyi karşılaştırır.
- Doğruluk Kontrolü : User'ın kimliğinin doğrulamasını yapar.
- Merkezi Yönetim : Şifre doğrulama mantığı tek bir fonksiyonda toplanır.

NOT: Şifre eşleşirse user'a erişim izni verilir.(Token veya Session oluşturulur.)
*/
func (u *User) CheckPassword(password string) bool {

	/*
		- "CompareHashAndPassword()", düz şifreyi tekrar hashleyip içerisinde ki "bcrypt" ile kendi doğrulama mantığını çalıştırır.
		- "u.PasswordHash" → veritabanında kayıtlı olan hashlenmiş şifre
		 -"password" → kullanıcının giriş ekranında yazdığı düz metin şifre

	*/
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash), // "u.PasswordHash" → veritabanında kayıtlı olan hashlenmiş şifre
		[]byte(password),       // "password" → kullanıcının giriş ekranında yazdığı düz metin şifre
	)

	return err == nil
}
