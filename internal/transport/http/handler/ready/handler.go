package ready

import (
	"net/http"
)

type Handler struct {
	readyUC ReadyUsecase
}

func NewHandler(readyUC ReadyUsecase) *Handler {
	return &Handler{
		readyUC: readyUC,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	if ok := h.readyUC.Execute(r.Context()); !ok {
		http.Error(w, "health check failed", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
