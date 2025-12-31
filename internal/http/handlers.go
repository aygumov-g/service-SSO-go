package http

import (
	"html/template"
	"net/http"
)

type handler struct {
	tmpl *template.Template
}

func newHandler(t *template.Template) *handler {
	return &handler{tmpl: t}
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Execute(w, map[string]string{
		"Title": "Главная",
	})
}
