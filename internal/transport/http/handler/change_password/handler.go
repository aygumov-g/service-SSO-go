package change_password

import (
	"encoding/json"
	"errors"
	"net/http"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Handler struct {
	change_passwordUC ChangePaswordUsecase
	identity          IdentityHTTP
}

func NewHandler(change_passwordUC ChangePaswordUsecase, identity IdentityHTTP) *Handler {
	return &Handler{
		change_passwordUC: change_passwordUC,
		identity:          identity,
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

	err := h.change_passwordUC.Execute(r.Context(), identity.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, account_d.ErrAccountNotFound):
			http.Error(w, "account not found", http.StatusConflict)
		case errors.Is(err, account_d.ErrSamePassword):
			http.Error(w, "new password must differ from old", http.StatusBadRequest)
		case errors.Is(err, account_d.ErrInvalidCredentials):
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
