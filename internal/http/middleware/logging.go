package middleware

import (
	"net/http"
	"time"

	"github.com/aygumov-g/service-SSO-go/internal/logger"
)

func Logging(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			log.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start).String())
		})
	}
}
