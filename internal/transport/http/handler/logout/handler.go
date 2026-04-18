package logout

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	logoutUC LogoutUsecase
}

func NewHandler(logoutUC LogoutUsecase) *Handler {
	return &Handler{
		logoutUC: logoutUC,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.post(w, r)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := h.logoutUC.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
