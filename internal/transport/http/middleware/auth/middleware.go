package auth

import (
	"net/http"
	"strings"
)

type Middleware struct {
	jwt      JWTService
	identity IdentityHTTP
}

func NewMiddleware(jwt JWTService, identity IdentityHTTP) *Middleware {
	return &Middleware{
		jwt:      jwt,
		identity: identity,
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

		ctx := m.identity.Upload(r.Context(), identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
