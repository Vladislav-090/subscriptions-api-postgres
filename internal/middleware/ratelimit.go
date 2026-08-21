package middleware

import (
	"net"
	"net/http"
	"subscriptions-api-postgres/internal/response"
)

type RateLimiter interface {
	Allow(key string) bool
}

func RateLimit(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				host = r.RemoteAddr
			}

			if !limiter.Allow(host) {
				response.WriteError(w, http.StatusTooManyRequests, "too many login attempts, try again later")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
