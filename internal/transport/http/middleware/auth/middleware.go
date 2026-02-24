package auth

import (
	"net/http"
	"strings"
)

type Middleware struct {
	jwt      JWTService
	identity IdentityHTTP
	accounts AccountService
}

func NewMiddleware(jwt JWTService, identity IdentityHTTP, accounts AccountService) *Middleware {
	return &Middleware{
		jwt:      jwt,
		identity: identity,
		accounts: accounts,
	}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")

		identity, err := m.jwt.Parse(token)
		if err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		if err = m.accounts.ValidateTokenVersion(r.Context(), identity); err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := m.identity.Upload(r.Context(), identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
