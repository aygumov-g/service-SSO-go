package http

import (
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl *template.Template
}

func NewHandler(t *template.Template) *Handler {
	return &Handler{tmpl: t}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, map[string]string{
		"Title": "Главная",
	})
}
