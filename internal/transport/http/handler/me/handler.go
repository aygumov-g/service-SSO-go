package me

import (
	"encoding/json"
	"errors"
	"net/http"

	srv_account "github.com/aygumov-g/service-SSO-go/internal/service/account"
)

type Handler struct {
	accounts AccountService
	identity IdentityHTTP
}

func NewHandler(accounts AccountService, identity IdentityHTTP) *Handler {
	return &Handler{
		accounts: accounts,
		identity: identity,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	identity, ok := h.identity.Unload(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	account, err := h.accounts.GetByID(r.Context(), identity.ID)
	if err != nil {
		switch {
		case errors.Is(err, srv_account.ErrAccountNotFound):
			http.Error(w, "account not found", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var resp accountResponse
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(account))
}
