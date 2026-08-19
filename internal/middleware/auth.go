package middleware

import (
	"context"
	"net/http"
	"strings"
	"subscriptions-api-postgres/internal/auth"
	"subscriptions-api-postgres/internal/models"
	"subscriptions-api-postgres/internal/response"
)

type TokenRevoker interface {
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
}

func Authenticate(jwtManager *auth.JWTManager, revoker TokenRevoker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || token == "" {
				response.WriteError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			claims, err := jwtManager.Parse(token)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			revoked, err := revoker.IsTokenRevoked(r.Context(), claims.ID)
			if err != nil {
				response.WriteServerError(w, "failed to check token status", err)
				return
			}
			if revoked {
				response.WriteError(w, http.StatusUnauthorized, "token has been revoked")
				return
			}

			ctx := auth.ContextWithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.ClaimsFromContext(r.Context())
			if !ok {
				response.WriteError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			if claims.Role != role {
				response.WriteError(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
