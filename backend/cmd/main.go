package main

import (
	"database/sql"
	"go-hexagonal/internal/adapters/handler"
	"go-hexagonal/internal/adapters/handler/middleware"
	"go-hexagonal/internal/adapters/limiter"
	"go-hexagonal/internal/adapters/repository"
	services "go-hexagonal/internal/core/service"
	"log"
	"net/http"
	"os"

	_ "go-hexagonal/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lib/pq" // Burada _ (underscore) olması ZORUNLUDUR! Go bu paketi sadece import eder ama sürücü kayıt (register) işlemini yapmaz. _ işareti, "paketin init() fonksiyonunu çalıştır ama başka fonksiyonlarını doğrudan kullanma" demektir ki sürücüler için gereken budur.
)

func main() {

	// 1. Database Connection Kurulumu
	connStr := os.Getenv("DB_URL") // DB_URL ortam değişkenini oku ("Docker-compose.yml" içinde ki "DB_URL" değişkenini kullanıyoruz.)

	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=wallet_db sslmode=disable"
		//log.Println("DB_URL bulunamadı, default ayarlar kullanılıyor.")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı kurulamadı: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Veritabanına ulaşılamadı: %v", err)
	}

	// 2. Repository Başlatma (Driven Adapter - Altyapı) başlatma işlemi
	walletRepo := repository.NewPostgreWalletRepository(db)
	auditRepo := repository.NewPostgreAuditRepository(db)

	// 3. Service (Core İş Mantığı) başlatma işlemi ve "repo" inject
	walletService := services.NewWalletService(walletRepo, auditRepo)

	// 4. Handler (Driving Adapter - API Katmanı) başlatma işlemi ve "service" inject
	walletHandler := handler.NewWalletHandler(walletService)

	// 5. Rate Limiter Start (SANİYEDE 2 İSTEK, 5 BURST KAPASİTE) - 6. Rate Limiter Create
	ratelimiter := limiter.NewInMemoryRateLimiter(2, 5)
	rateLimiterMiddleware := middleware.RateLimiterMiddleware(ratelimiter)

	// 7. HTTP ServeMux
	mux := http.NewServeMux()

	// 8. Statik Swagger Docs
	//fs := http.FileServer(http.Dir("./docs")) // statik dosyaları bir klasör listesi olarak listelemek.

	// 9. API Route
	//mux.HandleFunc("POST /wallets", walletHandler.Create)
	//mux.HandleFunc("GET /wallets/{id}", walletHandler.GetByID)
	//mux.HandleFunc("POST /wallets/{id}/deposit", walletHandler.Deposit)
	//mux.HandleFunc("POST /wallets/{id}/withdraw", walletHandler.Withdraw)
	//mux.HandleFunc("GET /wallets/{id}/transactions", walletHandler.GetTransactions)
	//mux.HandleFunc("GET /wallets/{id}/balance", walletHandler.GetBalance)
	//mux.HandleFunc("POST /wallets/{id}/transfer", walletHandler.Transfer)

	mux.Handle("POST /wallets",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.Create)))
	mux.Handle("GET /wallets/{id}",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.GetByID)))
	mux.Handle("POST /wallets/{id}/deposit",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.Deposit)))
	mux.Handle("POST /wallets/{id}/withdraw",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.Withdraw)))
	mux.Handle("GET /wallets/{id}/transactions",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.GetTransactionsByID)))
	mux.Handle("GET /wallets/{id}/balance",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.GetBalanceByID)))
	mux.Handle("POST /wallets/{id}/transfer",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.Transfer)))
	mux.Handle("POST /wallets/{id}/close",
		rateLimiterMiddleware(http.HandlerFunc(walletHandler.CloseWalletByID)))
	mux.Handle("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/swagger.json"),
	))

	// 10. HTTP Server Starting
	log.Println("Sunucu:8080 Portunda çalışıyor......")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
