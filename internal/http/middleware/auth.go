package middleware

import (
	"net/http"
	"strings"

	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
	"github.com/aygumov-g/service-SSO-go/internal/http/response"
)

func Auth(tokens auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.UnauthorizedError(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.UnauthorizedError(w, "invalid authorization header")
				return
			}

			userID, err := tokens.Parse(parts[1])
			if err != nil {
				response.UnauthorizedError(w, "invalid token")
				return
			}

			ctx := auth.ContextWithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
