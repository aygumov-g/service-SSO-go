package auth

import (
	"net/http"
	"strings"
)

type Middleware struct {
	get_meUC GetMeUsecase
	identity IdentityHTTP
	tokens   TokenProvider
}

func NewMiddleware(get_meUC GetMeUsecase, identity IdentityHTTP, tokens TokenProvider) *Middleware {
	return &Middleware{
		get_meUC: get_meUC,
		identity: identity,
		tokens:   tokens,
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

		identity, err := m.tokens.Parse(token)
		if err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		account, err := m.get_meUC.Execute(r.Context(), identity.ID)
		if err != nil || account.TokenVersion != identity.TokenVersion {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := m.identity.Upload(r.Context(), identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
