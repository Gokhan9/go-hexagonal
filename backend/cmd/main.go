package main

import (
	"go-hexagonal/internal/adapters/handler"
	"go-hexagonal/internal/adapters/handler/middleware"
	"go-hexagonal/internal/adapters/limiter"
	"go-hexagonal/internal/adapters/repository"
	services "go-hexagonal/internal/core/service"
	"log"
	"net/http"

	_ "go-hexagonal/docs"
)

func main() {

	// 1. Repository (Driven Adapter - Altyapı) başlatma işlemi
	repo := repository.NewMemoryWalletRepository()

	// 2. Service (Core İş Mantığı) başlatma işlemi ve "repo" inject
	walletService := services.NewWalletService(repo)

	// 3. Handler (Driving Adapter - API Katmanı) başlatma işlemi ve "service" inject
	walletHandler := handler.NewWalletHandler(walletService)

	// 4. Rate Limiter Start (SANİYEDE 2 İSTEK, 5 BURST KAPASİTE)
	ratelimiter := limiter.NewInMemoryRateLimiter(2, 5)

	// 5. Rate Limiter Create
	rateLimiterMiddleware := middleware.RateLimiterMiddleware(ratelimiter)

	// 6. HTTP ServeMux
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./docs")) // statik dosyaları bir klasör listesi olarak listelemek.

	// 7. API Route
	//mux.HandleFunc("POST /wallets", walletHandler.Create)
	mux.HandleFunc("GET /wallets/{id}", walletHandler.GetByID)
	mux.HandleFunc("POST /wallets/{id}/deposit", walletHandler.Deposit)
	mux.HandleFunc("POST /wallets/{id}/withdraw", walletHandler.Withdraw)
	mux.HandleFunc("GET /wallets/{id}/transactions", walletHandler.GetTransactions)
	mux.HandleFunc("GET /wallets/{id}/balance", walletHandler.GetBalance)
	mux.HandleFunc("POST /wallets/{id}/transfer", walletHandler.Transfer)
	//mux.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)

	mux.Handle("POST /wallets", rateLimiterMiddleware(http.HandlerFunc(walletHandler.Create)))
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", fs)) //statik dosyaları bir klasör listesi olarak listelemek.

	// 8. HTTP Server Starting
	log.Println("Sunucu:8080 Portunda çalışıyor......")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
