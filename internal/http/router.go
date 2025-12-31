package http

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(tmpl *template.Template, log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	h := NewHandler(tmpl)

	r.Use(RequestLogger(log))

	r.Handle("/stt/*", http.StripPrefix("/stt/", http.FileServer(http.Dir("web/statics"))))

	r.Get("/", h.Index)

	return r
}
