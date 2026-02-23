package change_password

import (
	"encoding/json"
	"errors"
	"net/http"

	account_srv "github.com/aygumov-g/service-SSO-go/internal/service/account"
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
	case http.MethodPost:
		h.post(w, r)
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	identity, ok := h.identity.Unload(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req change_passwordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := h.accounts.ChangePassword(r.Context(), identity.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, account_srv.ErrAccountNotFound):
			http.Error(w, "account not found", http.StatusConflict)
		case errors.Is(err, account_srv.ErrSamePassword):
			http.Error(w, "new password must differ from old", http.StatusBadRequest)
		case errors.Is(err, account_srv.ErrInvalidCredentials):
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
