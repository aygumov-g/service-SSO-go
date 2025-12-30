package http

import (
	"html/template"
	"log/slog"

	"github.com/aygumov-g/service-SSO-go/internal/http/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(tmpl *template.Template, log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	h := handlers.NewHandler(tmpl, log)

	r.Use(RequestLogger(log))

	r.Get("/", h.Index)

	return r
}
