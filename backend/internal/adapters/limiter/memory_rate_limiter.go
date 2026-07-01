package limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

type InMemoryRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex // farklı birçok goroutine'in(farklı IP'lerden gelen istekler) aynı anda map'e erişip limiter oluşturmaya çalışırsa sistem çöker, dolayısıyla mutex eklendi.
	ratelmt  rate.Limit
	burst    int
}

func NewInMemoryRateLimiter(ratelmt float64, burst int) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		ratelmt:  rate.Limit(ratelmt),
		burst:    burst,
	}
}

func (rl *InMemoryRateLimiter) Allow(key string) bool {

	rl.mu.Lock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.ratelmt, rl.burst) // ratelmt(saniye başına istek sayısı), burst(anlık patlama kapasitesi) system protection.
		rl.limiters[key] = limiter
	}
	rl.mu.Unlock()

	return limiter.Allow()
}
