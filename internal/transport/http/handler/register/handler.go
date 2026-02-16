package register

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
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := h.accounts.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, srv_account.ErrAccountAlreadyExists):
			http.Error(w, "account already exists", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}
