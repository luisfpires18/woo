package middleware

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const playerIDKey contextKey = "playerID"

// TokenValidator is a function that validates JWT tokens and returns the player ID.
type TokenValidator func(token string) (int64, error)

// Auth returns middleware that validates JWT Bearer tokens on protected routes.
func Auth(validate TokenValidator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			playerID, err := validate(parts[1])
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Store player ID in request context
			ctx := context.WithValue(r.Context(), playerIDKey, playerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// PlayerIDFromContext extracts the authenticated player ID from the request context.
func PlayerIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(playerIDKey).(int64)
	return id, ok
}
