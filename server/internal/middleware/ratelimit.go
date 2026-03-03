package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimit returns middleware that rate-limits requests per IP address.
func RateLimit(requestsPerSecond int) Middleware {
	var (
		mu       sync.Mutex
		limiters = make(map[string]*rate.Limiter)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			limiter, ok := limiters[ip]
			if !ok {
				limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond)
				limiters[ip] = limiter
			}
			mu.Unlock()

			if !limiter.Allow() {
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
