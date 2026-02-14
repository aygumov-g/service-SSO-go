package method

import "net/http"

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Handle(methods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range methods {
			if v == r.Method {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	})
}
