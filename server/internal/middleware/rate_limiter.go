package middleware

import (
	"net/http"
	"server/internal/response"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

func newRateLimiter() *limiter.Limiter {
	return tollbooth.NewLimiter(0.5, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tollbooth.LimitByRequest(newRateLimiter(), w, r)
		if err != nil {
			response.SendMessage(w, "rate limited exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
