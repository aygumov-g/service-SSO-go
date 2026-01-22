package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aygumov-g/service-SSO-go/internal/domain/auth"
)

type RegisterService interface {
	Register(ctx context.Context, login, password string) error
}

type RegisterHandler struct {
	auth RegisterService
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewRegisterHandler(auth RegisterService) *RegisterHandler {
	return &RegisterHandler{auth: auth}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := h.auth.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserAlreadyExists):
			http.Error(w, auth.ErrUserAlreadyExists.Error(), http.StatusConflict)
		default:
			http.Error(w, auth.ErrInternal.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct{}{})
}
