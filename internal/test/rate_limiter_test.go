package test

import (
	"go-hexagonal/internal/adapters/handler/middleware"
	"go-hexagonal/internal/adapters/limiter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiterMiddleware(t *testing.T) {

	// 1. Setup: Saniyede 1 istek, 1 burst (ani 1 ekstra istek hakkı) kapasiteli limiter. Kullanıcı aynı anda en fazla "1" hızlı istek yapabilir, sonra beklemek zorunda
	rl := limiter.NewInMemoryRateLimiter(1, 1)
	middleware.RateLimiterMiddleware(rl)

	// Basit "dummy" handler (middleware'i geçip geçmediğini kontrol et)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := middleware.RateLimiterMiddleware(rl)(nextHandler)

	// 2. Execution: Ardışık istekler

	// İlk istek (İzin verilmeli)
	req1 := httptest.NewRequest("POST", "/wallets", nil)
	req1.RemoteAddr = "192.168.1.1:1234"
	w1 := httptest.NewRecorder()      // istek oluşturulur
	handlerToTest.ServeHTTP(w1, req1) // gönderilir.

	// TEST ASSERTIONS: Başarılı Olursa "200", eğer farklıysa "errorf".
	// w1.Code, İLK İSTEĞİN STATUS CODE
	// w1.Code = 200 - test PASS
	// w1.Code ≠ 200 - test FAIL
	if w1.Code != http.StatusOK {
		t.Errorf("İlk istek başarısız oldu. Beklenen: 200, Gelen: %d", w1.Code)
	}

	// İkinci istek (Hemen ardından geldiği için limit aşılmalı)
	req2 := httptest.NewRequest("POST", "/wallets", nil)
	req2.RemoteAddr = "192.168.1.1:1234"
	w2 := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w2, req2)

	// TEST ASSERTIONS: Aynı IP'den 2. hızlı istek engellenmeli. 2. istek 429 değilse, "errorf"
	// w2.Code, İKİNCİ İSTEĞİN STATUS CODE
	// w2.Code = 429 - test PASS
	// w2.Code ≠ 429 - test FAIL
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("İkinci istek engellenmeli. Beklenen: 429, Gelen: %d", w2.Code)
	}
}
