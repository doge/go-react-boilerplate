package middleware

import (
	"net/http"
	"server/internal/response"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

func newRateLimiter() *limiter.Limiter {
	return tollbooth.NewLimiter(1, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
	})
}

var apiRateLimiter = newRateLimiter()

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tollbooth.LimitByRequest(apiRateLimiter, w, r)
		if err != nil {
			response.SendError(w, http.StatusTooManyRequests, "RATE_LIMITED", "Rate limit exceeded.")
			return
		}
		next.ServeHTTP(w, r)
	})
}
