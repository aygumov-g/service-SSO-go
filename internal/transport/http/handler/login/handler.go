package login

import (
	"encoding/json"
	"errors"
	"net/http"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Handler struct {
	loginUC LoginUsecase
}

func NewHandler(loginUC LoginUsecase) *Handler {
	return &Handler{
		loginUC: loginUC,
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

	tokens, err := h.loginUC.Execute(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, account_d.ErrInvalidCredentials) || errors.Is(err, account_d.ErrAccountNotFound):
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	var resp response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp.toResponse(tokens))
}
