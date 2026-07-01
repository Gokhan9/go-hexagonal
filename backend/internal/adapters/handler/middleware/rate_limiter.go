package middleware

import (
	"go-hexagonal/internal/core/ports"
	"net"
	"net/http"
)

func RateLimiterMiddleware(limiter ports.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// IP adresini al
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			if !limiter.Allow(ip) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Too many request", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
