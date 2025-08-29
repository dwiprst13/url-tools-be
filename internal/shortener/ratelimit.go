package shortener

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	tokens int
	last   time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
	rate    int
	burst   int
	intv    time.Duration
}

func NewRateLimiter(rate, burst int, intv time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*tokenBucket),
		rate:    rate,
		burst:   burst,
		intv:    intv,
	}
}

func (rl *RateLimiter) Allow(r *http.Request) bool {
	ip := clientIP(r)
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[ip]
	if !ok {
		rl.buckets[ip] = &tokenBucket{tokens: rl.burst - 1, last: now}
		return true
	}

	elapsed := now.Sub(b.last)
	refill := int(elapsed / rl.intv * time.Duration(rl.rate))
	if refill > 0 {
		b.tokens = min(rl.burst, b.tokens+refill)
		b.last = now
	}

	if b.tokens <= 0 {
		return false
	}

	b.tokens--
	return true
}

func clientIP(r *http.Request) string{
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	host, _, err :=net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}