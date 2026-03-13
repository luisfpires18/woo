package middleware

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit returns middleware that rate-limits requests per IP address.
// Stale entries are cleaned up periodically to prevent memory leaks.
// The cleanup goroutine stops when ctx is cancelled.
func RateLimit(ctx context.Context, requestsPerSecond int) Middleware {
	var (
		mu       sync.Mutex
		limiters = make(map[string]*rateLimiterEntry)
	)

	// Background cleanup goroutine: every 3 minutes, remove entries not seen in 5 minutes.
	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				for ip, entry := range limiters {
					if time.Since(entry.lastSeen) > 5*time.Minute {
						delete(limiters, ip)
					}
				}
				mu.Unlock()
			}
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr // fallback if no port present
			}

			mu.Lock()
			entry, ok := limiters[ip]
			if !ok {
				entry = &rateLimiterEntry{
					limiter:  rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
					lastSeen: time.Now(),
				}
				limiters[ip] = entry
			} else {
				entry.lastSeen = time.Now()
			}
			mu.Unlock()

			if !entry.limiter.Allow() {
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
