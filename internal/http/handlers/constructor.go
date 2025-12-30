package handlers

import (
	"html/template"
	"log/slog"
)

type Handler struct {
	tmpl *template.Template
	log  *slog.Logger
}

func NewHandler(t *template.Template, log *slog.Logger) *Handler {
	return &Handler{
		tmpl: t,
		log:  log,
	}
}
