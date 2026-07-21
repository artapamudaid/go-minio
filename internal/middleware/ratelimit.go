package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

// GlobalLimiter membatasi total trafik yang masuk ke Go Service
type GlobalLimiter struct {
	limiter *rate.Limiter
}

func NewGlobalLimiter(r rate.Limit, b int) *GlobalLimiter {
	return &GlobalLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

func GlobalRateLimit(g *GlobalLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !g.limiter.Allow() {
			http.Error(w, `{"success":false,"message":"Server Storage Sibuk (Global Rate Limit Exceeded)"}`, http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
