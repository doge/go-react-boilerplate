package middleware

import (
	"net/http"
	"net/url"
	"server/internal/models"
	"server/internal/response"
	"strings"
)

func TrustedOriginMiddleware(config *models.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := strings.TrimSpace(r.Header.Get("Origin"))
			if origin == "" {
				referer := strings.TrimSpace(r.Header.Get("Referer"))
				if referer != "" {
					parsed, err := url.Parse(referer)
					if err == nil && parsed.Scheme != "" && parsed.Host != "" {
						origin = parsed.Scheme + "://" + parsed.Host
					}
				}
			}

			if !config.IsAllowedOrigin(origin) {
				response.SendError(w, http.StatusForbidden, "CSRF_ORIGIN_DENIED", "Untrusted request origin.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
