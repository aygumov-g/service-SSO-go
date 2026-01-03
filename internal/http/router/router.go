package router

import "net/http"

type Router struct {
	handler http.Handler
}

func New() *Router {
	return &Router{
		handler: http.NewServeMux(),
	}
}

func (r *Router) Handle(patter string, handler http.Handler) {
	r.handler.(*http.ServeMux).Handle(patter, handler)
}

func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.handler = mw(r.handler)
}

func (r *Router) Handler() http.Handler {
	return r.handler
}
