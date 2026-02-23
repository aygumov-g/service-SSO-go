package refresh

import (
	"encoding/json"
	"errors"
	"net/http"

	session_srv "github.com/aygumov-g/service-SSO-go/internal/service/session"
)

type Handler struct {
	sessions SessionService
}

func NewHandler(sessions SessionService) *Handler {
	return &Handler{
		sessions: sessions,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.post(w, r)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	tokens, err := h.sessions.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, session_srv.ErrInvalidRefreshToken):
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	var resp refreshResponse
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(tokens))
}
