package refresh

import (
	"encoding/json"
	"errors"
	"net/http"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
	session_d "github.com/aygumov-g/service-SSO-go/internal/domain/session"
)

type Handler struct {
	refreshUC RefreshUsecase
}

func NewHandler(refreshUC RefreshUsecase) *Handler {
	return &Handler{
		refreshUC: refreshUC,
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

	tokens, err := h.refreshUC.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, session_d.ErrTokenNotFound):
			http.Error(w, "refresh token not found", http.StatusNotFound)
		case errors.Is(err, session_d.ErrTokenExpired):
			http.Error(w, "refresh token expired", http.StatusBadRequest)
		case errors.Is(err, session_d.ErrTokenRevoked):
			http.Error(w, "refresh token revoked", http.StatusBadRequest)
		case errors.Is(err, account_d.ErrAccountNotFound):
			http.Error(w, "account not found", http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	var resp response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(tokens))
}
