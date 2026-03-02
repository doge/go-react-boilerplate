package middleware

import (
	"context"
	"net/http"
	"server/internal/response"
	"server/internal/security"
	"strings"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func AuthMiddleware(tokenManager *security.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				response.SendError(w, http.StatusUnauthorized, "MISSING_AUTHORIZATION", "Missing authorization header.")
				return
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.SendError(w, http.StatusUnauthorized, "INVALID_AUTHORIZATION", "Invalid authorization header.")
				return
			}

			claims, err := tokenManager.ValidateAccessToken(parts[1])
			if err != nil {
				response.SendError(w, http.StatusUnauthorized, "INVALID_ACCESS_TOKEN", "Invalid access token.")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}
