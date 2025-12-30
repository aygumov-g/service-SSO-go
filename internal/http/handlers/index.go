package handlers

import (
	"net/http"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.log.Info("index page requested")

	h.tmpl.Execute(w, map[string]string{
		"Title": "Главная",
	})
}
