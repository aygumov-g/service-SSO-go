package register

import (
	"encoding/json"
	"errors"
	"net/http"

	account_d "github.com/aygumov-g/service-SSO-go/internal/domain/account"
)

type Handler struct {
	registerUC RegisterUsecase
}

func NewHandler(registerUC RegisterUsecase) *Handler {
	return &Handler{
		registerUC: registerUC,
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

	err := h.registerUC.Execute(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, account_d.ErrAccountAlreadyExists):
			http.Error(w, "account already exists", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}
