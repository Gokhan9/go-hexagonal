package ports

// Rate Limiter Port, sistem için hız sınırlayıcı yeteneği tanımlar. Driven Adapter'lar (Memory, Redis) bu interface'i implement eder...
type RateLimiter interface {
	Allow(key string) bool // İstek gelince, IP adresini (key) olarak alır, bu porttan haber bekler.
}
