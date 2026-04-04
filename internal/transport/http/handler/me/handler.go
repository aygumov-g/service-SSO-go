package me

import (
	"encoding/json"
	"errors"
	"net/http"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Handler struct {
	get_meUC GetMeUsecase
	identity IdentityHTTP
}

func NewHandler(get_meUC GetMeUsecase, identity IdentityHTTP) *Handler {
	return &Handler{
		get_meUC: get_meUC,
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

	account, err := h.get_meUC.Execute(r.Context(), identity.ID)
	if err != nil {
		switch {
		case errors.Is(err, account_d.ErrAccountNotFound):
			http.Error(w, "account not found", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	var resp response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(account))
}
