package handler

import (
	"html/template"
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/http/response"
)

type IndexData struct {
	Title string
}

type IndexHandler struct {
	tmpl *template.Template
}

func NewIndexHandler(tmpl *template.Template) *IndexHandler {
	return &IndexHandler{tmpl: tmpl}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	data := IndexData{
		Title: "Главная",
	}

	if err := h.tmpl.Execute(w, data); err != nil {
		response.InternalError(w, "template error")
	}
}
