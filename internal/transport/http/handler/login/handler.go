package login

import (
	"encoding/json"
	"errors"
	"net/http"

	srv_account "github.com/aygumov-g/service-SSO-go/internal/service/account"
)

type Handler struct {
	accounts AccountService
}

func NewHandler(accounts AccountService) *Handler {
	return &Handler{
		accounts: accounts,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.post(w, r)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	token, err := h.accounts.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, srv_account.ErrInvalidCredentials):
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	var resp loginResponse
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(token))
}
